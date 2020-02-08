let map = L.map('mapid').setView([51.505, -0.09], 13);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);



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

<<<<<<< HEAD
    let mapData = await response.json();
    
    for(let i=0; i<mapData.length; i++){
        
        let mapItem = mapData[i];
        let logo = "";
        if(mapItem.Vendor == "craigslist"){
            logo = "/static/img/craigslist_logo.png";
        }
        /*
        let myrep = await fetch('/static/pages/popup.html');
        let elem = "";
        if(myrep.ok){
            elem = await myrep.text();
        }
        
        elem.getElementById("Vendor").src = logo;
        elem.getElementById("Title").innerHTML = mapItem.Title;
        elem.getElementById("Price").innerHTML = mapItem.Price;
        elem.getElementById("Posted").innerHTML = mapItem.Posted;
        */
        let listing = L.popup()
        .setLatLng([mapItem.Lat, mapItem.Long])
        //.setContent(elem)
        .openOn(map);
        
    }

    let result = await response.json();

    console.log(result)
>>>>>>> 6dd3ad612d9c5ceb92ebc37e3d3f5cdcca8832b2
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
