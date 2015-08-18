package main

import (
       "os"
       "io/ioutil"
       "fmt"
       "os/user"
       "bufio"
       "time"
       "encoding/json"
       "math/rand"
       "bytes"
       "net/http"
       "log"
       "strings"
       "github.com/codegangsta/cli"
)

var domain  = "http://www.nitrc.org/ir"
var project = "PING"

func getSubject( reg string ) {
  cookie := getCurrentCookie()

  // first find out the experiment ID
  url := domain + "/data/archive/projects/" + project + "/subjects/" + reg + "/experiments"
  postData := []byte("{'format':'json'}")
  req, _ := http.NewRequest( "GET", url, bytes.NewReader(postData))
  req.Header.Set("Content-Type", "Content-Type: text/json; charset=utf-8")
  req.AddCookie(&cookie)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      panic(err)
  }   
  var f interface{}
  err = json.Unmarshal(body, &f)
  if err != nil {
      usr,_ := user.Current()
      f := usr.HomeDir + "/.cookieForSession"
      err = os.Remove(f)
      if err != nil {
        fmt.Printf("Error: could not delete cookie file %s\n", f)
      }
      fmt.Printf("Error: Login failed, please try to login again\n\n")
      fmt.Printf("%v\n%s\n\n", err, body)            
      return
  }
  m := f.(map[string]interface{})
  experiment := ""
  for _, v := range m {
     // this is ResultSet
     m2 := v.(map[string]interface{}) 
     for _, v2 := range m2 {
       switch vv := v2.(type) {
         case string:
           // total records
           // fmt.Printf("number of records is: %s\n", vv)
         default:
           // this should be Result, an array of definitions
           //for k3, v3 := range vv {
           vvv := vv.([]interface{})
           // only look at the first entry (object)
           vvvv := vvv[0].(map[string]interface{})
           for k4, v4 := range vvvv {
             if k4 == "ID" {
               //fmt.Printf("Found key: %s and value %s\n", k4, v4)
               experiment = v4.(string)               
             }
           }
           //}
       }
     }
  }
  if experiment == "" {
    fmt.Printf("Error: could not find an ID key for this subject")
    return
  }

  //url = "http://www.nitrc.org/ir/data/archive/projects/PING/subjects/" + reg + "/experiments/" + experiment + "/scans/3/resources/DICOM/files"
  url = domain + "/data/archive/projects/" + project + "/subjects/" + reg + "/experiments/" + experiment + "/resources/DICOM/files?format=zip"
  fmt.Printf("download all DICOM files for subject %s\n", reg)
  req2, _ := http.NewRequest( "GET", url, nil)
  req2.Header.Set("Content-Type", "Content-Type: application/zip")
  req2.AddCookie(&cookie)

  client2 := &http.Client{}
  resp2, err := client2.Do(req2)
  if err != nil {
      panic(err)
  }
  defer resp2.Body.Close()

  filename := fmt.Sprintf("%s.zip", reg)
  fp, _ := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
  defer fp.Close() // defer close

  buf := make([]byte, 131072)
  var all []byte
  for {
    n, _ := resp2.Body.Read(buf)
    if n == 0 {
      break
    }
    all = append(all, buf[:n]...)
    fmt.Printf("\033[2K[%3.2fmb] %s\r", float64(len(all))/1024.0/1024.0, fp.Name())
  }
  fmt.Printf("\nwrite to disk...")
  _, err = fp.Write(all[:len(all)])
  if err != nil {
    log.Println(err)
  }
  fmt.Printf("\033[2K\r")  
}

func getListOfSubjects( reg string ) {
  cookie := getCurrentCookie()

  //url := "http://www.nitrc.org/ir/data/archive/projects/PING/subjects/?format=json"
  url := domain + "/data/archive/projects/" + project + "/subjects?columns=label,gender,age&format=json"
  req, _ := http.NewRequest( "GET", url, nil /* bytes.NewReader(postData) */)
  req.Header.Set("Content-Type", "Content-Type: text/json; charset=utf-8")
  req.AddCookie(&cookie)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      panic(err)
  }
  var f interface{}
  err = json.Unmarshal(body, &f)
  if err != nil {
      //fmt.Printf("%v\n%s\n\n", err, body)
      usr,_ := user.Current()
      f := usr.HomeDir + "/.cookieForSession"
      err = os.Remove(f)
      if err != nil {
        fmt.Printf("Error: could not delete cookie file %s\n", f)
      }
      fmt.Printf("Error: Login failed, please try to login again\n\n")
      fmt.Printf("%v\n%s\n\n", err, body)      
      return
  }
  m := f.(map[string]interface{})
  //count := 0
  var subjects [][3]string
  for _, v := range m {
     // this is ResultSet
     m2 := v.(map[string]interface{}) 
     for _, v2 := range m2 {
       switch vv := v2.(type) {
         case string:
           // total records
           // fmt.Printf("number of records is: %s\n", vv)
         default:
           // this should be Result, an array of definitions
           //for k3, v3 := range vv {
           vvv := vv.([]interface{})
           // only look at the first entry (object)
           for _, v3 := range vvv {
             vvvv := v3.(map[string]interface{})
             var subj string
             var gender string
             var age string
             for k4, v4 := range vvvv {
               if k4 == "label" {
                 // fmt.Printf("Found subject:  %s and value %s\n", k4, v4)
                 // subjects = append(subjects, v4.(string))
                 subj = v4.(string)
               }
               if k4 == "gender" {
                 gender = v4.(string)                 
               }
               if k4 == "age" {
                 age = v4.(string)                 
               }
             }
             // add information to output
             subjects = append(subjects, [3]string{ subj, age, gender})
           }
       }
     }
  }
  if len(subjects) == 0 {
    fmt.Printf("Error: could not find any subjects for PING")
    return
  }

  fmt.Printf("SubjID, Age, Gender\n")
  fmt.Printf("# subjects found: %d\n", len(subjects))
  for _, v := range subjects {
     fmt.Printf("%s, %s, %s\n", v[0], v[1], v[2]) 
  }
    
}

// the expiration time stored might not be good, ignore for now
type Settings struct {
  CookieName string
  CookieValue string
  CookieExpires time.Time
}

func getNewCookie() ( http.Cookie ) {
  // ask the user for user name and password for NITRC
  fmt.Printf("User name: ")
  reader := bufio.NewReader(os.Stdin)
  username, _ := reader.ReadString('\n')
  username = strings.Replace(username, "\n", "", -1)
  
  fmt.Printf("Password: ")
  var pass []byte = make( []byte, 100 )
  os.Stdin.Read( pass )
  println()

  // todo: we should do an insecure connection :-/
  url := domain + "/data/JSESSION"
  req, _ := http.NewRequest( "POST", url, nil)
  req.Proto = "HTTP/1.1"
  req.ProtoMinor = 1
  req.Header.Set("Content-Type", "Content-Type: text/plain;charset=ISO-8859-1")
  req.SetBasicAuth(username, string(pass))
  
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  _, _ = ioutil.ReadAll(resp.Body)
  
  c := resp.Cookies()
  var cookieCode = c[0]
  
  // save the cookie to our file for next time, we don't want to enter the user name and password all the time
  saveCurrentCookie( *cookieCode )
  return *cookieCode
}

func saveCurrentCookie( cookie http.Cookie ) {
  // fmt.Printf("Save this cookie: %s", cookie)
  usr,_ := user.Current()
  dir := usr.HomeDir + "/.cookieForSession"
  // if file exists load it here
  var s Settings
  s.CookieName = cookie.Name
  s.CookieValue = cookie.Value
  s.CookieExpires = cookie.Expires

  b, err := json.MarshalIndent(s, "", "  ")
  if err != nil { panic(err) }

  fi, err := os.Create(dir)
  if err != nil { panic(err) }
  defer func() {
    if err := fi.Close(); err != nil {
      panic(err)
    }
  }()
  n := len(b)
  if _, err := fi.Write(b[:n]); err != nil {
    panic(err)
  }
}

func getCurrentCookie() ( http.Cookie ) {
  // try to read the current cookie file for this project
  var m Settings
  usr,_ := user.Current()
  f := usr.HomeDir + "/.cookieForSession"
  buf, err := ioutil.ReadFile(f)
  if err != nil {
    c := getNewCookie()
    m.CookieName = c.Name
    m.CookieValue = c.Value
    m.CookieExpires = c.Expires
  } else {
    err = json.Unmarshal(buf[:len(buf)], &m)
    if err != nil {
      panic(err)
    }  
  }
  return http.Cookie{Name: m.CookieName, Value: m.CookieValue, Expires: m.CookieExpires}
}

func main() {
     rand.Seed(time.Now().UTC().UnixNano())

     app := cli.NewApp()
     app.Name    = "pingdata"
     app.Usage   = "Download " + project + " data from " + domain + ".\n\n" +
                   "This program uses the XNAT REST API to download subject data. Start by listing\n" +
                   "the subjects available:\n\n" +
                   " > pingdata list\n\n" +
                   "To download the data of a specific subject call:\n\n" +
                   " > pingdata pull PXXXXX\n" +
                   "where PXXXXX is the subject identification number." 
     app.Version = "0.0.1"
     app.Author  = "Hauke Bartsch"
     app.Email   = "HaukeBartsch@gmail.com"

     app.Commands = []cli.Command{
     {
        Name:      "pull",
        ShortName: "p",
        Usage:     "Retrieve subject data as zip",
        Description: "Download subject data as a zip file into the current directory.\n\n",
        Action: func(c *cli.Context) {
          if len(c.Args()) < 1 {
            fmt.Printf("Error: indiscriminate downloads are not supported, supply a subject name\n")
          } else {
            getSubject( c.Args().First()) // only deliver output
          }
        },
     },
     {
        Name:      "list",
        ShortName: "l",
        Usage:     "Retrieve list of subjects",
        Description: "Display a list of subjects with Subject ID, age, and gender as a comma-separated-value (csv) list.\n\n",
        Action: func(c *cli.Context) {
          if len(c.Args()) == 1 {
            getListOfSubjects( c.Args()[0] )
          } else {
            getListOfSubjects( ".*")
          }
        },
     },
  }

  app.Run(os.Args)
}