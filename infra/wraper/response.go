package wraper

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Response struct {
	Total   any           `json:"total,omitempty" xml:"total,omitempty"`
	Data    any           `json:"data,omitempty" xml:"data,omitempty"`
	Code    int           `json:"code" xml:"code"`
	Msg     string        `json:"msg,omitempty" xml:"msg,omitempty"`
	Blob    BlobDef       `json:"-" xml:"-"`
	HTML    string        `json:"-" xml:"-"`
	Status  int           `json:"-" xmd:"-"`
	Mime    MimeType      `json:"-" xml:"-"`
	Request *http.Request `json:"-" xml:"-"`
}

func Error(msg interface{}) *Response {
	str := ""
	if msg == nil {
		return &Response{
			Code:   200,
			Status: http.StatusOK,
		}
	}
	switch msg := msg.(type) {
	case string:
		str = msg
	case error:
		str = msg.Error()
	default:
		str = ""
	}
	return &Response{
		Msg:    str,
		Code:   404,
		Mime:   MineJson,
		Status: http.StatusOK,
	}
}
func Json(input interface{}) *Response {
	return &Response{
		Data:   input,
		Code:   200,
		Mime:   MineJson,
		Status: http.StatusOK,
	}
}
func Empty() *Response {
	return &Response{
		Code:   200,
		Mime:   MineJson,
		Status: http.StatusOK,
	}
}
func OkData(input interface{}) *Response {
	return &Response{
		Data:   input,
		Code:   200,
		Mime:   MineJson,
		Status: http.StatusOK,
	}
}
func OkMsg(msg string) *Response {
	return &Response{
		Data:   nil,
		Code:   200,
		Msg:    msg,
		Mime:   MineJson,
		Status: http.StatusOK,
	}
}
func Blob(input BlobDef) *Response {
	return &Response{
		Blob:   input,
		Code:   200,
		Mime:   MineBlob,
		Status: http.StatusOK,
	}
}

func HTML(input string) *Response {
	return &Response{
		HTML:   input,
		Code:   200,
		Mime:   MineHtml,
		Status: http.StatusOK,
	}
}

func XML(input any) *Response {
	return &Response{
		Data:   input,
		Code:   200,
		Mime:   MineXml,
		Status: http.StatusOK,
	}
}

// 获得msg
func (w *Response) HttpStatus(status int) *Response {
	w.Status = status
	return w
}

// 获得msg
func (w *Response) WithMine(mime MimeType) *Response {
	w.Mime = mime
	return w
}

// 获得msg
func (w *Response) WithMsg(msg string) *Response {
	w.Msg = msg
	return w
}

// 获得msg
func (Response *Response) EncodeJSON(writer http.ResponseWriter) (err error) {

	return json.NewEncoder(writer).Encode(writer)

}

// 获得msg
func (Response *Response) Encode(writer http.ResponseWriter) (err error) {

	mime := Response.Mime
	err = nil
	if mime == MineHtml || mime == MineText {
		writer.WriteHeader(Response.Status)
		_, err = fmt.Fprintf(writer, Response.HTML)
	} else if mime == MineJson {
		writer.WriteHeader(Response.Status)
		err = json.NewEncoder(writer).Encode(Response)
	} else if mime == MineXml {
		writer.WriteHeader(Response.Status)
		err = xml.NewEncoder(writer).Encode(Response)
	} else if mime == MineBlob {
		contenttype := "application/octet-stream"
		if Response.Blob.ContentType != "" {
			contenttype = Response.Blob.ContentType
		}
		writer.Header().Set("Content-Type", contenttype)
		if Response.Blob.Name != "" {
			writer.Header().Set("Content-Disposition", "attachment; filename="+Response.Blob.Name) // 用来指定下载下来的文件名
		}
		writer.WriteHeader(Response.Status)
		filebytes := Response.Blob.File
		_, err = writer.Write(filebytes)

	} else {
		writer.WriteHeader(Response.Status)
		err = json.NewEncoder(writer).Encode(Response)
	}
	return err
}

// 获得msg
func (w *Response) WithTotal(total any) *Response {
	w.Total = total
	return w
}

// 获得msg
func (w *Response) WithHTML(html string) *Response {
	w.HTML = html
	return w
}

// 获得msg
func (w *Response) WithBlob(blob BlobDef) *Response {
	w.Blob = blob
	return w
}

// 获得msg
func (w *Response) WithError(err error) *Response {
	if err != nil {
		w.Code = 400
		w.Msg = err.Error()
	}
	return w
}
