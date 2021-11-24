import { useState } from "react"


const Input = ({ onShorten }) => {

    const [url,setUrl] = useState("")

    const onSubmit = (e) => {
        e.preventDefault()

        if (!url || url.length < 5) {
            alert("Please enter a URL!")
            return
        }

        console.log("On Submit")
        onShorten({url})
        setUrl("")
        
    }


    return (
        <form className="input_form" onSubmit={onSubmit}>
            <div className="input_div">
                <input className="input_field" type="text" placeholder="https://www.google.com/" value={url} onChange={(e) => setUrl(e.target.value)}/>
            </div>
            <input className="btn btn-block" type="submit" value="Shorten URL"/>
        </form>
    )
}

export default Input
