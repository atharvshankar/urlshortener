function getShortUrl(){
    var longurl = document.getElementById('longurl').value;
    if(longurl == ""){
        alert("Please enter a URL!")
        document.getElementById('longurl').value="";
    }
    // axios.get(longurl).then((res) => {
    //     if(res.status == 404){
    //         alert("Please enter valid url");
    //     }
    //     alert("Thanks for valid url");
    // }).catch((err) => {
    //     alert("Please enter valid url");
    // })
    axios.post("http://localhost:3000/shorten/",{
        "url": longurl
    })
    .then((res) => {
        console.log(res.data['long_url']);
        var link = document.getElementById("placeholder-shorturl");
        link.setAttribute("href", res.data['short_url']);
        document.getElementById("placeholder-shorturl").text = res.data['short_url']
    })
    return
}