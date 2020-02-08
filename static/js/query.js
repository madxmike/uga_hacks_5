let map = L.map('mapid').setView([51.505, -0.09], 13);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


var marker = L.marker([51.5, -0.09]).addTo(map);
search.onsubmit = async (e) => {
    e.preventDefault();
    console.log("test")
    let form = new FormData(search);
    let priceMin = form.get("price_min")
    let priceMax = form.get("price_max")
    console.log(priceMin)
    console.log(priceMax)
    if (priceMin === "") {
        form.set("price_min", "0")
    }
    if (priceMax === "") {
        form.set("price_max", "1000")
    }
    form.append("bounds", map.getBounds().toBBoxString());
    console.log(JSON.stringify(Object.fromEntries(form)))
    let response = await fetch(window.location.origin + '/api/search', {
      method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(Object.fromEntries(form))
    });

    let result = await response.json();

    console.log(result)
  };

document.getElementsByClassName("currency").onblur =function (){

    //number-format the user input
    this.value = parseFloat(this.value.replace(/,/g, ""))
        .toFixed(2)
        .toString()
        .replace(/\B(?=(\d{3})+(?!\d))/g, ",");

    //set the numeric value to a number input
    document.getElementById("number").value = this.value.replace(/,/g, "")

}