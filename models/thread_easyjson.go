// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2d00218DecodeParkBaseParkDbModels(in *jlexer.Lexer, out *Threads) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Threads, 0, 8)
			} else {
				*out = Threads{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *Thread
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(Thread)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeParkBaseParkDbModels(out *jwriter.Writer, in Threads) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Threads) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeParkBaseParkDbModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Threads) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeParkBaseParkDbModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Threads) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeParkBaseParkDbModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Threads) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeParkBaseParkDbModels(l, v)
}
func easyjson2d00218DecodeParkBaseParkDbModels1(in *jlexer.Lexer, out *ThreadUpdate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		case "title":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeParkBaseParkDbModels1(out *jwriter.Writer, in ThreadUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeParkBaseParkDbModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeParkBaseParkDbModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeParkBaseParkDbModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeParkBaseParkDbModels1(l, v)
}
func easyjson2d00218DecodeParkBaseParkDbModels2(in *jlexer.Lexer, out *Thread) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "id":
			out.ID = int64(in.Int64())
		case "message":
			out.Message = string(in.String())
		case "slug":
			if in.IsNull() {
				in.Skip()
				out.Slug = nil
			} else {
				if out.Slug == nil {
					out.Slug = new(string)
				}
				*out.Slug = string(in.String())
			}
		case "title":
			out.Title = string(in.String())
		case "votes":
			out.Votes = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeParkBaseParkDbModels2(out *jwriter.Writer, in Thread) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Created).MarshalJSON())
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Forum))
	}
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	if in.Slug != nil {
		const prefix string = ",\"slug\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.Slug))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.Votes != 0 {
		const prefix string = ",\"votes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Votes))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Thread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeParkBaseParkDbModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeParkBaseParkDbModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeParkBaseParkDbModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeParkBaseParkDbModels2(l, v)
}