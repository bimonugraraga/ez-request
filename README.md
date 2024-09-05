# ez-request

## How To Use

### EzRequest
This method is for sending reguler http request. You can also set ```json StatusCodeConstraint``` to limit status code that you want.
``` go
temp := RequestParams{
    Ctx:                  context.Background(), // mandatory
	Method:               http.MethodPost, // mandatory
	URL:                  "", // mandatory
    TimeoutMs:            1000,
	StatusCodeConstraint: []int{200},
}
res, err := temp.EzRequest()
defer res.Body.Close()

body, err := io.ReadAll(res.Body)
type randomStruct struct {}

err = json.Unmarshal(res.Body, &randomStruct)
```

### EzRetriableRequest() 
This method is for sending http request with retriable mechanism. You can set how many retry ```json Attemps```, interval between retry ```json BackoffMs```, and timeout for each request ```json TimeoutMs```
``` go
temp := RequestParams{
    Ctx:                  context.Background(), // mandatory
	Method:               http.MethodPost, // mandatory
	URL:                  "", // mandatory
	Attempts:             2,
	BackoffMs:            1000,
	TimeoutMs:            1000,
	StatusCodeConstraint: []int{200},
}
res, err := temp.EzRetriableRequest() 
defer res.Body.Close()

body, err := io.ReadAll(res.Body)
type randomStruct struct {}

err = json.Unmarshal(res.Body, &randomStruct)
```

### EzDoIt and EzDoItRetriable
EzDoIt and EzDoItRetriable let you make your custom http client and request. You can still set the retriable params and status code constraints
``` go
client := &http.Client{}
req, err := http.NewRequestWithContext(context.Background(),http.MethodPost, "https://...", bytes.NewBuffer([]byte("")))
temp := RequestParams{
	Attempts:             2,
	BackoffMs:            1000,
	TimeoutMs:            1000,
	StatusCodeConstraint: []int{200},
}
res, err := temp.EzDoit(req, client)
res, err := temp.EzDoItRetriable(req, client)
defer res.Body.Close()

body, err := io.ReadAll(res.Body)
type randomStruct struct {}

err = json.Unmarshal(res.Body, &randomStruct)
```