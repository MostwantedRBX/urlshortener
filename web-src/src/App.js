import { useState } from "react"
import Header from "./components/Header";
import Input from "./components/Input";

function App() {
  const [shortenedUrl,setShortenedUrl] = useState("")

  async function shortenUrl(url) {
    const res = await fetch("https://srtlink.net/put/",{ // Prod
      // const res = await fetch("http://localhost:8080/put/",{ // Testing  
      method:"POST",
      headers:{
        'Content-Type':"application/json",
      },
      body: JSON.stringify(url)
    }) 

    const data = await res.json()
    setShortenedUrl(data.url)
  }
  
  function copyButton() {
    
    navigator.clipboard.writeText(shortenedUrl);
    
    alert("Copied " + shortenedUrl + " to clipboard.")
  }

  return (
    <div className="container">
      <Header />
      {console.log("Rendered")}
      <Input onShorten={shortenUrl} />
      <p className="urlP">
        {shortenedUrl ? shortenedUrl:"Shorten a link!"}
      </p>
      
      {shortenedUrl ? <button className="btn btn-block" onClick={copyButton}>Copy to Clipboard</button>:null}
      {shortenedUrl ? <button className="btn btn-block" onClick={async () => { 
        // await fetch("http://localhost:8080/qr/",{ // Testing
          await fetch("https://srtlink.net//qr/",{ // Prod
          method:"POST",
          headers:{
            'Content-Type':"application/json",
          },
          body: JSON.stringify({url:shortenedUrl})
        })
        .then(res => res.blob())
        .then(blob => URL.createObjectURL(blob))
        .then(url => {
          const a = document.createElement('a');
          a.href = url
          a.download = 'qr-code.png'
          a.click()
        })
      }}>
      Download QR Code
      </button>:null}
    </div>
  );
}

export default App;
