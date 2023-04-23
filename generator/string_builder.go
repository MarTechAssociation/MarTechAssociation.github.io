// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import "strings"

type StringBuilder struct {
	arrayBegin          bool
	justAppendSeparator bool // prevent separator to append duplicatedly
	sb                  *strings.Builder
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{
		sb: &strings.Builder{},
	}
}

func (s *StringBuilder) write(str string) {
	s.arrayBegin = false
	s.justAppendSeparator = false
	s.sb.WriteString(str)
}

func (s *StringBuilder) Len() int {
	return s.sb.Len()
}

func (s *StringBuilder) BeginArray() {
	s.write("[")
	s.arrayBegin = true
}

func (s *StringBuilder) EndArray() {
	s.write("]")
}

func (s *StringBuilder) AppendIfNotFirstItemInArray(str string) {
	if s.arrayBegin || s.justAppendSeparator {
		return
	}
	s.write(str)
	s.justAppendSeparator = true
}

func (s *StringBuilder) Append(str string) {
	s.write(str)
}

func (s *StringBuilder) String() string {
	return s.sb.String()
}
