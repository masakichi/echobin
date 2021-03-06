package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const maxByteCount = 100 << 10
const maxDelay = 10 // seconds

// @Summary  The request's query parameters.
// @Tags     HTTP methods
// @Produce  json
// @Success  200  {object}  getMethodResponse
// @Router   /get [get]
func getMethodHandler(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, &getMethodResponse{
		Args:    getArgs(c),
		Headers: getHeaders(c),
		Origin:  getOrigin(c),
		URL:     getURL(c),
	}, "  ")
}

// @Summary  The request's query parameters.
// @Tags     HTTP methods
// @Accept   json
// @Accept   mpfd
// @Accept   x-www-form-urlencoded
// @Produce  json
// @Success  200  {object}  otherMethodResponse
// @Router   /post [post]
// @Router   /put [put]
// @Router   /patch [patch]
// @Router   /delete [delete]
func otherMethodHandler(c echo.Context) error {
	data := ""
	files := getFiles(c)
	form := getForm(c)
	if len(files) == 0 && len(form) == 0 {
		data = getData(c)
	}
	res := otherMethodResponse{}
	res.Args = getArgs(c)
	res.Data = data
	res.Files = files
	res.Form = form
	res.Headers = getHeaders(c)
	res.JSON = getJSON(c)
	res.Origin = getOrigin(c)
	res.URL = getURL(c)
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

// @Summary  Returns the requester's IP Address.
// @Tags     Request inspection
// @Produce  json
// @Success  200  {object}  requestIPResponse  "The Requester’s IP Address."
// @Router   /ip [get]
func requestIPHandler(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, &requestIPResponse{
		Origin: getOrigin(c),
	}, "  ")
}

// @Summary  Return the incoming request's HTTP headers.
// @Tags     Request inspection
// @Produce  json
// @Success  200  {object}  requestHeadersResponse  "The request’s headers."
// @Router   /headers [get]
func requestHeadersHandler(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, &requestHeadersResponse{
		Headers: getHeaders(c),
	}, "  ")
}

// @Summary  Return the incoming requests's User-Agent header.
// @Tags     Request inspection
// @Produce  json
// @Success  200  {object}  requestUserAgentResponse  "The request’s User-Agent header."
// @Router   /user-agent [get]
func requestUserAgentHandler(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, &requestUserAgentResponse{
		UserAgent: getUserAgent(c),
	}, "  ")
}

type weightedCode struct {
	weight float64
	code   int
}

func chooseStatusCode(weightedCodes []weightedCode) int {
	var code int
	var total float64
	var cumWeights []float64
	for _, wc := range weightedCodes {
		total += wc.weight
		cumWeights = append(cumWeights, total)
	}
	rand.Seed(time.Now().UnixNano())
	x := rand.Float64() * total
	for i, cumWeight := range cumWeights {
		if cumWeight > x {
			code = weightedCodes[i].code
			break
		}
	}
	return code
}

// @Summary   Return status code or random status code if more than one are given
// @Tags      Status codes
// @Produce   plain
// @Param     codes  path  string  true  "codes"
// @Response  100    "Informational responses"
// @Response  200    "Success"
// @Response  300    "Redirection"
// @Response  400    "Client Errors"
// @Response  500    "Server Errors"
// @Router    /status/{codes} [delete]
// @Router    /status/{codes} [get]
// @Router    /status/{codes} [patch]
// @Router    /status/{codes} [post]
// @Router    /status/{codes} [put]
func statusCodesHandler(c echo.Context) error {
	codes, _ := url.PathUnescape(c.Param("codes"))
	if !strings.Contains(codes, ",") {
		code, err := strconv.Atoi(codes)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid status code")
		}
		return c.NoContent(code)
	}

	var weightedCodes []weightedCode
	var _code, _weight string
	for _, choice := range strings.Split(codes, ",") {
		if !strings.Contains(choice, ":") {
			_code = choice
			_weight = "1"
		} else {
			s := strings.SplitN(choice, ":", 2)
			_code, _weight = s[0], s[1]
		}
		code, err := strconv.Atoi(_code)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid status code")
		}
		weight, err := strconv.ParseFloat(_weight, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid status code")
		}
		weightedCodes = append(weightedCodes, weightedCode{weight, code})
	}
	return c.NoContent(chooseStatusCode(weightedCodes))
}

//go:embed static/moby.html
var sampleHTML []byte

// @Summary   Returns a simple HTML document.
// @Tags      Response formats
// @Produce   html
// @Response  200  "An HTML page."
// @Router    /html [get]
func serveHTMLHandler(c echo.Context) error {
	return c.HTMLBlob(http.StatusOK, sampleHTML)
}

//go:embed static/sample.xml
var sampleXML []byte

// @Summary   Returns a simple XML document.
// @Tags      Response formats
// @Produce   xml
// @Response  200  "An XML document."
// @Router    /xml [get]
func serveXMLHandler(c echo.Context) error {
	return c.XMLBlob(http.StatusOK, sampleXML)
}

//go:embed static/sample.json
var sampleJSON []byte

// @Summary   Returns a simple JSON document.
// @Tags      Response formats
// @Produce   json
// @Response  200  "An JSON document."
// @Router    /json [get]
func serveJSONHandler(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, sampleJSON)
}

const ROBOTS_TXT = `User-agent: *
Disallow: /deny
`

// @Summary   Returns some robots.txt rules.
// @Tags      Response formats
// @Produce   plain
// @Response  200  "Robots file"
// @Router    /robots.txt [get]
func serveRobotsTXTHandler(c echo.Context) error {
	return c.String(http.StatusOK, ROBOTS_TXT)
}

//go:embed static/deny.txt
var denyTXT string

// @Summary   Returns page denied by robots.txt rules.
// @Tags      Response formats
// @Produce   plain
// @Response  200  "Denied message"
// @Router    /deny [get]
func serveDenyHandler(c echo.Context) error {
	return c.String(http.StatusOK, denyTXT)
}

//go:embed static/sample-utf8.html
var sampleUTF8HTML []byte

// @Summary   Returns a UTF-8 encoded body.
// @Tags      Response formats
// @Produce   html
// @Response  200  "Encoded UTF-8 content."
// @Router    /encoding/utf8 [get]
func serveUTF8HTMLHandler(c echo.Context) error {
	return c.Blob(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, sampleUTF8HTML)
}

// @Summary   Returns GZip-encoded data.
// @Tags      Response formats
// @Produce   json
// @Response  200              "GZip-encoded data."
// @Param     accept-encoding  header  string  false  "Accept-Encoding"  default(gzip)
// @Router    /gzip [get]
func serveGzipHandler(c echo.Context) error {
	res := gzippedResponse{}
	res.Origin = getOrigin(c)
	res.Headers = getHeaders(c)
	res.Method = c.Request().Method
	if strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "gzip") {
		res.Gzipped = true
	}
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

// @Summary   Returns Deflate-encoded data.
// @Tags      Response formats
// @Produce   json
// @Response  200              "Defalte-encoded data."
// @Param     accept-encoding  header  string  false  "Accept-Encoding"  default(deflate)
// @Router    /deflate [get]
func serveDeflateHandler(c echo.Context) error {
	res := deflatedResponse{}
	res.Origin = getOrigin(c)
	res.Headers = getHeaders(c)
	res.Method = c.Request().Method
	if strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "deflate") {
		res.Deflated = true
	}
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

// @Summary   Returns Brotli-encoded data.
// @Tags      Response formats
// @Produce   json
// @Response  200              "Brotli-encoded data."
// @Param     accept-encoding  header  string  false  "Accept-Encoding"  default(br)
// @Router    /brotli [get]
func serveBrotliHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	w := brotli.HTTPCompressor(c.Response().Writer, c.Request())
	defer w.Close()

	res := brotliResponse{}
	res.Origin = getOrigin(c)
	res.Headers = getHeaders(c)
	res.Method = c.Request().Method
	if strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "br") {
		res.Brotli = true
	}

	// TODO: better to make a middleware for brotli
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(&res)
}

// @Summary   Decodes base64url-encoded string.
// @Tags      Dynamic data
// @Produce   plain
// @Param     value  path  string  true  "Encoded base64 content"  default(RUNIT0JJTiBpcyBhd2Vzb21l)
// @Response  200    "Decoded base64 content."
// @Router    /base64/{value} [get]
func base64Handler(c echo.Context) error {
	value, _ := url.PathUnescape(c.Param("value"))
	value = strings.TrimSpace(value)
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		bytes, err = base64.URLEncoding.DecodeString(value)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Incorrect Base64 data try: RUNIT0JJTiBpcyBhd2Vzb21l")
		}
	}
	return c.String(http.StatusOK, string(bytes))
}

// @Summary   Returns n random bytes generated with given seed
// @Tags      Dynamic data
// @Produce   octet-stream
// @Param     n     path   int  true   "number of bytes"
// @Param     seed  query  int  false  "seed"
// @Response  200   "Bytes."
// @Router    /bytes/{n} [get]
func generateBytesHandler(c echo.Context) error {
	n := c.Param("n")
	seed := c.QueryParam("seed")
	intN, err := strconv.Atoi(n)
	if err != nil || intN < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of bytes")
	}
	if intN > maxByteCount {
		intN = maxByteCount
	}
	seedInt, err := strconv.Atoi(seed)
	if err == nil {
		rand.Seed(int64(seedInt))
	}
	bytes := make([]byte, intN)
	rand.Read(bytes)
	return c.Blob(http.StatusOK, echo.MIMEOctetStream, bytes)
}

// @Summary   Returns a delayed response (max of 10 seconds).
// @Tags      Dynamic data
// @Produce   json
// @Param     delay  path  int  true  "delay"
// @Response  200    "A delayed response."
// @Router    /delay/{delay} [delete]
// @Router    /delay/{delay} [get]
// @Router    /delay/{delay} [patch]
// @Router    /delay/{delay} [post]
// @Router    /delay/{delay} [put]
func delayHandler(c echo.Context) error {
	delay := c.Param("delay")
	intDelay, err := strconv.Atoi(delay) // TODO: support float type delay
	if err != nil || intDelay < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of delay")
	}
	if intDelay > maxDelay {
		intDelay = maxDelay
	}
	time.Sleep(time.Duration(intDelay) * time.Second)
	data := ""
	files := getFiles(c)
	form := getForm(c)
	if len(files) == 0 && len(form) == 0 {
		data = getData(c)
	}
	return c.JSONPretty(http.StatusOK, &delayResponse{
		Args:    getArgs(c),
		Data:    data,
		Files:   files,
		Form:    form,
		Headers: getHeaders(c),
		Origin:  getOrigin(c),
		URL:     getURL(c),
	}, "  ")
}

type dripParams struct {
	// The amount of time (in seconds) over which to drip each byte
	Duration float64 `query:"duration" default:"2"`
	// The number of bytes to respond with
	Numbytes int `query:"numbytes" default:"10"`
	// The response code that will be returned
	Code int `query:"code" default:"200"`
	// The amount of time (in seconds) to delay before responding
	Delay float64 `query:"delay" default:"2"`
}

// @Summary   Drips data over a duration after an optional initial delay.
// @Tags      Dynamic data
// @Produce   octet-stream
// @Param     dripParams  query  dripParams  true  "dripParams"
// @Response  200         "A dripped response."
// @Router    /drip [get]
func dripHandler(c echo.Context) error {
	dp := &dripParams{
		Duration: 2,
		Numbytes: 10,
		Code:     200,
		Delay:    2,
	}
	if err := c.Bind(dp); err != nil {
		return err
	}

	if dp.Delay < 0 {
		dp.Delay = 0
	} else if dp.Delay > 10 {
		dp.Delay = 10
	}

	if dp.Duration < 0.1 {
		dp.Duration = 0.1 // Minimum duration = 100 Millisecond
	} else if dp.Duration > 60 {
		dp.Duration = 60
	}

	if dp.Numbytes < 0 {
		dp.Numbytes = 0
	} else if dp.Numbytes > 10<<20 {
		dp.Numbytes = 10 << 20 // Millisecond
	}

	time.Sleep(time.Duration(dp.Delay*1000) * time.Millisecond)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	c.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(dp.Numbytes))
	c.Response().WriteHeader(dp.Code) // TODO: validate status code?

	remainBytes := dp.Numbytes
	times := int(dp.Duration / 0.1)
	chunkLength := dp.Numbytes
	if times > 1 {
		chunkLength = dp.Numbytes/times + 1
	}
	if chunkLength == 1 {
		pause := int(dp.Duration*1000) / remainBytes
		for remainBytes > 0 {
			if _, err := c.Response().Write([]byte{'*'}); err != nil {
				return err
			}
			c.Response().Flush()
			time.Sleep(time.Duration(pause) * time.Millisecond)
			remainBytes--
		}
	} else {
		for remainBytes > 0 {
			var length int
			if remainBytes > chunkLength {
				length = chunkLength
			} else {
				length = remainBytes
			}
			if _, err := c.Response().Write(bytes.Repeat([]byte{'*'}, length)); err != nil {
				return err
			}
			c.Response().Flush()
			time.Sleep(100 * time.Millisecond)
			remainBytes -= length
		}
	}
	return nil
}

// TODO: swaggo doesn't support struct when type is path...
type linksParams struct {
	// The amount of links
	N int `param:"n" default:"1"`
	// Offset starts from 0
	Offset int `param:"offset" default:"0"`
}

//go:embed templates/links.html
var linksTemplate string

// @Summary   Generate a page containing n links to other pages which do the same.
// @Tags      Dynamic data
// @Produce   html
// @Param     n       path  int  true  "The amount of links"   default(1)
// @Param     offset  path  int  true  "Offset starts from 0"  default(0)
// @Response  200     "HTML links."
// @Router    /links/{n}/{offset} [get]
func linksHandler(c echo.Context) error {
	lp := &linksParams{
		N:      1,
		Offset: 0,
	}

	if err := c.Bind(lp); err != nil {
		return err
	}

	if lp.N < 1 {
		lp.N = 1 // Minimum 1
	} else if lp.N > 200 {
		lp.N = 200 // Maximum 200
	}

	t := template.Must(template.New("links").Parse(linksTemplate))
	seq := make([]int, lp.N, lp.N)

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, map[string]interface{}{
		"seq":     seq,
		"offset":  lp.Offset,
		"reverse": c.Echo().Reverse,
	}); err != nil {
		return err
	}
	// Seems echo will not set content-length if lenght is over 2048
	// see also: https://github.com/labstack/echo/pull/366
	return c.Blob(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, buf.Bytes())
}

type rangeParams struct {
	Numbytes  int     `param:"numbytes"`
	ChunkSize int     `query:"chunk_size"`
	Duration  float64 `query:"duration"`
}

// @Summary   Streams n random bytes generated with given seed, at given chunk size per packet.
// @Tags      Dynamic data
// @Produce   octet-stream
// @Param     numbytes    path   int     true   "The amount of bytes"  default(10)
// @Param     chunk_size  query  int     false  "chunk_size"
// @Param     duration    query  number  false  "duration"
// @Response  200         "Bytes"
// @Router    /range/{numbytes} [get]
func rangeHandler(c echo.Context) error {
	rp := &rangeParams{
		ChunkSize: 10 << 10,
	}
	if err := c.Bind(rp); err != nil {
		return err
	}
	if rp.Numbytes <= 0 || rp.Numbytes > maxByteCount {
		c.Response().Header().Set("ETag", fmt.Sprintf("range%d", rp.Numbytes))
		c.Response().Header().Set("Accept-Ranges", "bytes")
		return echo.NewHTTPError(http.StatusNotFound, "number of bytes must be in the range (0, 102400]")
	}
	if rp.ChunkSize < 1 {
		rp.ChunkSize = 1
	}
	if rp.Duration < 0 {
		rp.Duration = 0
	} else if rp.Duration > 60 {
		rp.Duration = 60
	}
	first, last := getRequestRange(c.Request().Header.Get("Range"), rp.Numbytes)
	if first > last || last >= rp.Numbytes {
		c.Response().Header().Set("ETag", fmt.Sprintf("range%d", rp.Numbytes))
		c.Response().Header().Set("Accept-Ranges", "bytes")
		c.Response().Header().Set("Content-Range", fmt.Sprintf("bytes */%d", rp.Numbytes))
		c.Response().Header().Set(echo.HeaderContentLength, "0")
		return c.NoContent(http.StatusRequestedRangeNotSatisfiable)
	}

	contentRange := fmt.Sprintf("bytes %d-%d/%d", first, last, rp.Numbytes)
	contentLength := fmt.Sprintf("%d", last-first+1)
	code := http.StatusPartialContent
	if first == 0 && last == rp.Numbytes-1 {
		code = http.StatusOK
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	c.Response().Header().Set("ETag", fmt.Sprintf("range%d", rp.Numbytes))
	c.Response().Header().Set("Accept-Ranges", "bytes")
	c.Response().Header().Set(echo.HeaderContentLength, contentLength)
	c.Response().Header().Set("Content-Range", contentRange)
	c.Response().WriteHeader(code)

	pausePerByte := rp.Duration * 1000 / float64(last-first+1) // Millisecond
	cursor := first
	for cursor <= last {
		chunk := rp.ChunkSize
		if chunk >= last-cursor {
			chunk = last - cursor + 1
		}
		pause := pausePerByte * float64(chunk)
		time.Sleep(time.Duration(pause) * time.Millisecond)
		bytes := make([]byte, chunk)
		for i := cursor; i < cursor+chunk; i++ {
			bytes[i-cursor] = byte('a' + i%26)
		}
		if _, err := c.Response().Write(bytes); err != nil {
			return err
		}
		c.Response().Flush()
		cursor += chunk
	}
	return nil
}

type streamBytesParams struct {
	N         int `param:"n"`
	Seed      int `query:"seed"`
	ChunkSize int `query:"chunk_size"`
}

// @Summary   Streams n random bytes generated with given seed, at given chunk size per packet.
// @Tags      Dynamic data
// @Produce   octet-stream
// @Param     n           path   int  true   "The amount of bytes"
// @Param     seed        query  int  false  "seed"
// @Param     chunk_size  query  int  false  "chunk_size"
// @Response  200         "Bytes"
// @Router    /stream-bytes/{n} [get]
func streamBytesHandler(c echo.Context) error {
	sbp := &streamBytesParams{
		ChunkSize: 10 << 10,
	}
	if err := c.Bind(sbp); err != nil {
		return err
	}

	if sbp.N > maxByteCount {
		sbp.N = maxByteCount
	}
	if sbp.ChunkSize < 1 {
		sbp.ChunkSize = 1
	}
	if c.QueryParams().Has("seed") {
		rand.Seed(int64(sbp.Seed))
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	c.Response().WriteHeader(http.StatusOK)
	remainBytes := sbp.N
	for remainBytes > 0 {
		chunk := sbp.ChunkSize
		if sbp.ChunkSize > remainBytes {
			chunk = remainBytes
		}
		bytes := make([]byte, chunk)
		rand.Read(bytes)
		if _, err := c.Response().Write(bytes); err != nil {
			return err
		}
		c.Response().Flush()
		remainBytes -= chunk
	}
	return nil
}

// @Summary   Stream n JSON responses
// @Tags      Dynamic data
// @Produce   json
// @Param     n    path  int  true  "The amount of JSON objects"  default(10)
// @Response  200  "Streamed JSON responses."
// @Router    /stream/{n} [get]
func streamHandler(c echo.Context) error {
	n := c.Param("n")
	intN, err := strconv.Atoi(n)
	if err != nil || intN < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of JSON objects")
	}
	if intN > 100 {
		intN = 100
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	res := streamResponse{
		Args:    getArgs(c),
		Headers: getHeaders(c),
		Origin:  getOrigin(c),
		URL:     getURL(c),
	}
	enc := json.NewEncoder(c.Response())
	for i := 0; i < intN; i++ {
		res.ID = i
		if err := enc.Encode(res); err != nil {
			return err
		}
		c.Response().Flush()
	}
	return nil
}

// @Summary   Return a UUID4.
// @Tags      Dynamic data
// @Produce   json
// @Response  200  "A UUID4."
// @Router    /uuid [get]
func UUIDHandler(c echo.Context) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, &UUIDResponse{
		UUID: uuid,
	}, "  ")
}

//go:embed static/images/sample.webp
var sampleWebP []byte

//go:embed static/images/sample.svg
var sampleSVG []byte

//go:embed static/images/sample.jpeg
var sampleJPEG []byte

//go:embed static/images/sample.png
var samplePNG []byte

// @Summary   Returns a simple image of the type suggest by the Accept header.
// @Tags      Images
// @Produce   image/webp
// @Produce   image/svg+xml
// @Produce   image/jpeg
// @Produce   image/png
// @Produce   image/*
// @Response  200  "An image."
// @Router    /image [get]
func imageHandler(c echo.Context) error {
	switch accept := strings.ToLower(c.Request().Header.Get(echo.HeaderAccept)); {
	case strings.Contains(accept, "image/webp"):
		return imageWebPHandler(c)
	case strings.Contains(accept, "image/svg+xml"):
		return imageSVGHandler(c)
	case strings.Contains(accept, "image/jpeg"):
		return imageJPEGHandler(c)
	case strings.Contains(accept, "image/png"), strings.Contains(accept, "image/*"):
		return imagePNGHandler(c)
	default:
		return echo.NewHTTPError(http.StatusNotAcceptable, "Client did not request a supported media type.")
	}
}

// @Summary   Returns a simple WEBP image.
// @Tags      Images
// @Produce   image/webp
// @Response  200  "A WEBP image."
// @Router    /image/webp [get]
func imageWebPHandler(c echo.Context) error {
	return c.Blob(http.StatusOK, "image/webp", sampleWebP)
}

// @Summary   Returns a simple SVG image.
// @Tags      Images
// @Produce   image/svg+xml
// @Response  200  "An SVG image."
// @Router    /image/svg [get]
func imageSVGHandler(c echo.Context) error {
	return c.Blob(http.StatusOK, "image/svg+xml", sampleSVG)
}

// @Summary   Returns a simple JPEG image.
// @Tags      Images
// @Produce   image/jpeg
// @Response  200  "A JPEG image."
// @Router    /image/jpeg [get]
func imageJPEGHandler(c echo.Context) error {
	return c.Blob(http.StatusOK, "image/jpeg", sampleJPEG)
}

// @Summary   Returns a simple PNG image.
// @Tags      Images
// @Produce   image/png
// @Response  200  "A PNG image."
// @Router    /image/png [get]
func imagePNGHandler(c echo.Context) error {
	return c.Blob(http.StatusOK, "image/png", samplePNG)
}

// @Summary   Returns cookie data.
// @Tags      Cookies
// @Produce   json
// @Response  200  "Cookies"
// @Router    /cookies [get]
func getCookiesHandler(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, &cookiesResponse{
		Cookies: getCookies(c),
	}, "  ")
}

// @Summary   Deletes cookie(s) as provided by the query string and redirects to cookie list.
// @Tags      Cookies
// @Produce   json
// @Param     freeform  query  string  false  "freeform"
// @Response  200       "Redirect to cookie list"
// @Router    /cookies/delete [get]
func deleteCookiesHandler(c echo.Context) error {
	for k := range c.QueryParams() {
		cookie, _ := c.Cookie(k)
		if cookie != nil {
			cookie.MaxAge = -1
			cookie.Path = "/"
			c.SetCookie(cookie)
		}
	}
	return c.Redirect(http.StatusFound, c.Echo().URI(getCookiesHandler))
}

// @Summary   Sets cookie(s) as provided by the query string and redirects to cookie list.
// @Tags      Cookies
// @Produce   json
// @Param     freeform  query  string  false  "freeform"
// @Response  200       "Redirect to cookie list"
// @Router    /cookies/set [get]
func setCookiesInQueryHandler(c echo.Context) error {
	for k, v := range c.QueryParams() {
		cookie := &http.Cookie{
			Name:  k,
			Value: v[0],
			Path:  "/",
		}
		c.SetCookie(cookie)
	}
	return c.Redirect(http.StatusFound, c.Echo().URI(getCookiesHandler))
}

// @Summary   Sets a cookie and redirects to cookie list.
// @Tags      Cookies
// @Produce   json
// @Param     name   path  string  true  "name"
// @Param     value  path  string  true  "value"
// @Response  200    "Set cookies and redirects to cookie list."
// @Router    /cookies/set/{name}/{value} [get]
func setCookiesInPathHandler(c echo.Context) error {
	name := c.Param("name")
	value := c.Param("value")
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, c.Echo().URI(getCookiesHandler))
}

type redirectToParams struct {
	URL        string `query:"url" form:"url"`
	StatusCode int    `query:"status_code" form:"status_code"`
}

func redirectToHandler(c echo.Context) error {
	rp := &redirectToParams{
		StatusCode: http.StatusFound,
	}
	if err := c.Bind(rp); err != nil {
		return err
	}
	if rp.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "url is required")
	}
	if rp.StatusCode < 300 || rp.StatusCode > 308 {
		rp.StatusCode = http.StatusFound
	}
	return c.Redirect(rp.StatusCode, rp.URL)
}

// @Summary   302/3XX Redirects to the given URL.
// @Tags      Redirects
// @Produce   plain
// @Param     url          query  string  true   "url"
// @Param     status_code  query  int     false  "status_code"
// @Response  302          "A redirection."
// @Router    /redirect-to [get]
func getRedirectToHandler(c echo.Context) error {
	return redirectToHandler(c)
}

// @Summary   302/3XX Redirects to the given URL.
// @Tags      Redirects
// @Accept    x-www-form-urlencoded
// @Produce   plain
// @Param     url          formData  string  true   "url"
// @Param     status_code  formData  int     false  "status_code"
// @Response  302          "A redirection."
// @Router    /redirect-to [delete]
// @Router    /redirect-to [patch]
// @Router    /redirect-to [post]
// @Router    /redirect-to [put]
func otherRedirectToHandler(c echo.Context) error {
	return redirectToHandler(c)
}

func redirect(absolute bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		n := c.Param("n")
		intN, err := strconv.Atoi(n)
		if err != nil || intN < 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid number of redirection times")
		}
		var redirectURI string
		if intN == 1 {
			redirectURI = c.Echo().URI(getMethodHandler)
		} else {
			h := relativeRedirectHandler
			if absolute {
				h = absoluteRedirectHandler
			}
			redirectURI = c.Echo().URI(h, strconv.Itoa(intN-1))
		}
		if absolute {
			redirectURI = c.Scheme() + "://" + c.Request().Host + redirectURI
		}
		return c.Redirect(http.StatusFound, redirectURI)
	}
}

// @Summary   302 Redirects n times.
// @Tags      Redirects
// @Produce   plain
// @Param     n         path   int     true   "n"
// @Param     absolute  query  string  false  "absolute"  default(false)
// @Response  302       "A redirection."
// @Router    /redirect/{n} [get]
func redirectHandler(c echo.Context) error {
	absolute := c.QueryParam("absolute") == "true"
	return redirect(absolute)(c)
}

// @Summary   Relatively 302 Redirects n times.
// @Tags      Redirects
// @Produce   plain
// @Param     n    path  int  true  "n"
// @Response  302  "A redirection."
// @Router    /relative-redirect/{n} [get]
func relativeRedirectHandler(c echo.Context) error {
	return redirect(false)(c)
}

// @Summary   Absolutely 302 Redirects n times.
// @Tags      Redirects
// @Produce   plain
// @Param     n    path  int  true  "n"
// @Response  302  "A redirection."
// @Router    /absolute-redirect/{n} [get]
func absoluteRedirectHandler(c echo.Context) error {
	return redirect(true)(c)
}

// @Summary   Returns anything passed in request data.
// @Tags      Anything
// @Accept    json
// @Accept    mpfd
// @Accept    x-www-form-urlencoded
// @Produce   json
// @Response  200  "Anything passed in request"
// @Router    /anything [delete]
// @Router    /anything [get]
// @Router    /anything [patch]
// @Router    /anything [post]
// @Router    /anything [put]
// @Router    /anything/{anything} [delete]
// @Router    /anything/{anything} [get]
// @Router    /anything/{anything} [patch]
// @Router    /anything/{anything} [post]
// @Router    /anything/{anything} [put]
func anythingHandler(c echo.Context) error {
	data := ""
	files := getFiles(c)
	form := getForm(c)
	if len(files) == 0 && len(form) == 0 {
		data = getData(c)
	}
	res := anythingResponse{}
	res.Args = getArgs(c)
	res.Data = data
	res.Files = files
	res.Form = form
	res.Headers = getHeaders(c)
	res.JSON = getJSON(c)
	res.Origin = getOrigin(c)
	res.URL = getURL(c)
	res.Method = c.Request().Method
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

// @Summary   Returns a 304 if an If-Modified-Since header or If-None-Match is present. Returns the same as a GET otherwise.
// @Tags      Response inspection
// @Produce   json
// @Response  200                "Normal response"
// @Response  304                "Not modified"
// @Param     If-Modified-Since  header  string  false  "If-Modified-Since"
// @Param     If-None-Match      header  string  false  "If-None-Match"
// @Router    /cache [get]
func cacheHandler(c echo.Context) error {
	header := c.Request().Header
	if header.Get(echo.HeaderIfModifiedSince) == "" && header.Get("If-None-Match") == "" {
		etag := uuid.New()
		c.Response().Header().Set(echo.HeaderLastModified, time.Now().UTC().Format(http.TimeFormat))
		c.Response().Header().Set("ETag", fmt.Sprintf("%x", etag[:]))
		return getMethodHandler(c)
	}
	// TODO: seems we need to return etag in header as well.
	// see also: https://www.ietf.org/rfc/rfc2616.txt
	return c.NoContent(http.StatusNotModified)
}

// @Summary   Sets a Cache-Control header for n seconds.
// @Tags      Response inspection
// @Produce   json
// @Response  200    "Cache control set"
// @Param     value  path  int  true  "Seconds"
// @Router    /cache/{value} [get]
func cacheDurationHandler(c echo.Context) error {
	value := c.Param("value")
	maxAge, err := strconv.Atoi(value)
	if err != nil || maxAge < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of seconds")
	}
	c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	return getMethodHandler(c)
}

// @Summary   Assumes the resource has the given etag and responds to If-None-Match and If-Match headers appropriately.
// @Tags      Response inspection
// @Produce   json
// @Response  200            "Normal response"
// @Response  304            "Not modified"
// @Response  412            "Precondition failed"
// @Param     etag           path    string  true   "etag"
// @Param     If-Match       header  string  false  "If-Match"
// @Param     If-None-Match  header  string  false  "If-None-Match"
// @Router    /etag/{etag} [get]
func etagHandler(c echo.Context) error {
	etag := c.Param("etag")
	ifNoneMatch := c.Request().Header.Get("If-None-Match")
	ifMatch := c.Request().Header.Get("If-Match")
	if ifNoneMatch != "" {
		if strings.Contains(ifNoneMatch, etag) || strings.Contains(ifNoneMatch, "*") {
			c.Response().Header().Set("ETag", etag)
			return c.NoContent(http.StatusNotModified)
		}
	} else if ifMatch != "" {
		if !strings.Contains(ifMatch, etag) && !strings.Contains(ifMatch, "*") {
			return c.NoContent(http.StatusPreconditionFailed)
		}
	}
	c.Response().Header().Set("ETag", etag)
	return getMethodHandler(c)
}

// @Summary   Returns a set of response headers from the query string.
// @Tags      Response inspection
// @Produce   json
// @Response  200       "Normal response"
// @Param     freeform  query  string  false  "freeform"
// @Router    /response-headers [get]
// @Router    /response-headers [post]
func responseHeadersHandler(c echo.Context) error {
	contentLength := 0
	body := map[string]interface{}{}
	body[echo.HeaderContentType] = echo.MIMEApplicationJSONCharsetUTF8
	for k, v := range c.QueryParams() {
		if len(v) == 1 {
			body[k] = v[0]
		} else {
			body[k] = v
		}
		c.Response().Header()[k] = v
	}
	var bs []byte
	// this for loop is going to write content-length into body
	for {
		body[echo.HeaderContentLength] = fmt.Sprintf("%d", contentLength)
		bs, _ = json.MarshalIndent(body, "", "  ")
		if len(bs) == contentLength {
			break
		}
		contentLength = len(bs)
	}
	return c.JSONBlob(http.StatusOK, bs)
}

// @Summary   Prompts the user for authorization using HTTP Basic Auth.
// @Tags      Auth
// @Produce   json
// @Response  200     "Sucessful authentication."
// @Response  401     "Unsuccessful authentication."
// @Param     user    path  string  true  "user"
// @Param     passwd  path  string  true  "passwd"
// @Router    /basic-auth/{user}/{passwd} [get]
func basicAuthHandler(c echo.Context) error {
	res := map[string]interface{}{
		"authenticated": true,
		"user":          c.Param("user"),
	}
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

// @Summary   Prompts the user for authorization using bearer authentication.
// @Tags      Auth
// @Produce   json
// @Response  200            "Sucessful authentication."
// @Response  401            "Unsuccessful authentication."
// @Param     Authorization  header  string  false  "Authorization"
// @Router    /bearer [get]
func bearerHandler(c echo.Context) error {
	authorization := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
	if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
		c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Bearer")
		return c.NoContent(http.StatusUnauthorized)
	}
	token := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")
	res := map[string]interface{}{
		"authenticated": true,
		"token":         token,
	}
	return c.JSONPretty(http.StatusOK, &res, "  ")
}

//go:embed static/swagger-ui
var swaggerUIFiles embed.FS

func swaggerUIHandler(c echo.Context) error {
	swaggerUIRoot, _ := fs.Sub(swaggerUIFiles, "static/swagger-ui")
	assetHandler := http.FileServer(http.FS(swaggerUIRoot))
	return echo.WrapHandler(assetHandler)(c)
}

//go:embed docs/swagger.json
var swaggerDoc []byte

func swaggerDocHandler(c echo.Context) error {
	doc := make(map[string]interface{})
	json.Unmarshal(swaggerDoc, &doc)
	docInfo, _ := doc["info"].(map[string]interface{})
	if c.Scheme() == "https" {
		doc["schemes"] = []string{"https"}
	} else {
		doc["schemes"] = []string{"http", "https"}
	}
	docInfo["version"] = fmt.Sprintf("%s-%s", version, revision)
	return c.JSON(http.StatusOK, doc)
}

//go:embed static/form.html
var formHTML []byte

func formHandler(c echo.Context) error {
	return c.HTMLBlob(http.StatusOK, formHTML)
}
