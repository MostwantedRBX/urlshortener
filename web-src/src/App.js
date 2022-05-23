import { useState } from "react"
import Header from "./components/Header";
import Input from "./components/Input";

function App() {
  const [shortenedUrl,setShortenedUrl] = useState("")

  async function shortenUrl(url) {
    console.log(url.url)
    console.log("On shortenUrl")
    const res = await fetch("http://167.172.240.248:8080/put/",{
      method:"POST",
      headers:{
        'Content-Type':"application/json",
      },
      body: JSON.stringify(url)
    }) 

    const data = await res.json()
    console.log(data.url)
    setShortenedUrl(data.url)
  }

  return (
    <div className="container">
      <Header />
      {console.log("Rendered")}
      <Input onShorten={shortenUrl} />
      <p className="urlP">
        {shortenedUrl ? shortenedUrl:""}
      </p>
    </div>
  );
}

export default App;
