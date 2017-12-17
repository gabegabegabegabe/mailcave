package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/tambchop/mailcave/logging"
)

const (
	mimeMsgPostKey  = "mime_message"
	contentTypeKey  = "Content-Type"
	jsonContentType = "application/json"
)

// ErrorContent holds the content of an HTTP response that specifies an error that occurred.
type ErrorContent struct {
	Err string `json:"error"`
}

// StoreSuccessContent holds the content of a successful archive store operation.
type StoreSuccessContent struct {
	// ID is the database ID of the stored email message.
	ID string `json:"id"`
}

// Archivist services mail archival requests.
type Archivist struct {
	archive Archive
	logger  *logging.Logger
}

// NewArchivist creates a Archivist.
func NewArchivist(archive Archive, logger *logging.Logger) *Archivist {
	return &Archivist{
		archive: archive,
		logger:  logger,
	}
}

// Start causes the archivist to start running and archiving mail when requested.
func (a *Archivist) Start(ctx context.Context, addr string) error {

	err := a.archive.Open()
	if err != nil {
		return fmt.Errorf("failed to open archive: %s", err)
	}
	defer a.archive.Close()

	router := httprouter.New()
	router.GET("/", handleIndex)
	router.POST("/store", a.handleStore)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		return err
	}

	return nil
}

func handleIndex(respWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprint(respWriter, "the cave is too dark to have a UI presently =/")
}

func (a *Archivist) handleStore(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	mimeStr := req.Form.Get(mimeMsgPostKey)

	mimeMsg, err := mail.ReadMessage(strings.NewReader(mimeStr))
	if err != nil {
		a.sendErrorResponse(rw, "400 Bad Request", 400, fmt.Errorf("bad or missing mime_message: %s", err), req)
		return
	}

	msg, err := MsgFromMIMEMessage(mimeMsg)
	if err != nil {
		a.sendErrorResponse(rw, "400 Bad Request", 400, err, req)
		return
	}

	id, err := a.archive.ArchiveMessage(msg)
	if err != nil {
		a.sendErrorResponse(rw, "500 Internal Server Error", 500, err, req)
		return
	}

	a.sendStoreSuccessResponse(rw, id, req)
}

func (a *Archivist) sendErrorResponse(rw http.ResponseWriter, status string, code int, errMsg error, req *http.Request) {

	ec := &ErrorContent{errMsg.Error()}
	body := ""
	jsonBytes, err := json.Marshal(ec)
	if err != nil {
		a.logger.Printf("very strange... failed to marshal ErrorContent with error %s", err)
		body = err.Error()
	} else {
		body = string(jsonBytes)
	}

	rw.Header().Set(contentTypeKey, jsonContentType)
	rw.WriteHeader(code)
	fmt.Fprintf(rw, "%s", body)
}

func (a *Archivist) sendStoreSuccessResponse(rw http.ResponseWriter, id string, req *http.Request) {

	ssc := &StoreSuccessContent{id}
	body := ""
	jsonBytes, err := json.Marshal(ssc)
	if err != nil {
		a.logger.Printf("very strange... failed to marshal StoreSuccessContentw with error %s", err)
		body = err.Error()
	} else {
		body = string(jsonBytes)
	}

	rw.Header().Set(contentTypeKey, jsonContentType)
	rw.WriteHeader(201)
	fmt.Fprintf(rw, "%s", body)
}
