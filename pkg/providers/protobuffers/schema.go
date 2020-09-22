package protobuffers

import (
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jhump/protoreflect/desc"
)

// NewSchema constructs a new schema manifest from the given file descriptors
func NewSchema(descriptors []*desc.FileDescriptor) specs.Schemas {
	result := make(specs.Schemas, 0)

	for _, descriptor := range descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result[message.GetFullyQualifiedName()] = NewMessage("", message)
		}
	}

	return result
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(path string, descriptor *desc.MessageDescriptor) *specs.Property {
	fields := descriptor.GetFields()
	result := &specs.Property{
		Path:        path,
		Name:        descriptor.GetFullyQualifiedName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    1,
		Template: specs.Template{
			Message: &specs.Message{
				Properties: make(map[string]*specs.Property, len(fields)),
			},
		},
		Options: specs.Options{},
	}

	for _, field := range fields {
		result.Message.Properties[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
	}

	return result
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(path string, descriptor *desc.FieldDescriptor) *specs.Property {
	result := &specs.Property{
		Path:        path,
		Name:        descriptor.GetName(),
		Description: descriptor.GetSourceInfo().GetLeadingComments(),
		Position:    descriptor.GetNumber(),
		Options:     specs.Options{},
	}

	if descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_ENUM {
		enum := descriptor.GetEnumType()
		keys := map[string]*specs.EnumValue{}
		positions := map[int32]*specs.EnumValue{}

		for _, value := range enum.GetValues() {
			result := &specs.EnumValue{
				Key:         value.GetName(),
				Position:    value.GetNumber(),
				Description: value.GetSourceInfo().GetLeadingComments(),
			}

			keys[value.GetName()] = result
			positions[value.GetNumber()] = result
		}

		result.Enum = &specs.Enum{
			Name:        enum.GetName(),
			Description: enum.GetSourceInfo().GetLeadingComments(),
			Keys:        keys,
			Positions:   positions,
		}

		return result
	}

	if descriptor.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE {
		fields := descriptor.GetMessageType().GetFields()
		result.Message = &specs.Message{
			Properties: make(map[string]*specs.Property, len(fields)),
		}

		for _, field := range fields {
			result.Message.Properties[field.GetName()] = NewProperty(template.JoinPath(path, field.GetName()), field)
		}

		return result
	}

	result.Scalar = &specs.Scalar{
		Type:  Types[descriptor.GetType()],
		Label: Labels[descriptor.GetLabel()],
	}

	return result
}
