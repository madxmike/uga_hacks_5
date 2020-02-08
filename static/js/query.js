let map = L.map('mapid').setView([33.7490, -84.3880], 11);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);
let markers = L.layerGroup();


search.onsubmit = async (e) => {
    e.preventDefault();
    document.getElementById("spin").style.display = "inline-block";
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

    let mapData = await response.json();
    
    markers.clearLayers(); 
    for(let i=0; i<mapData.length; i++){
        
        let mapItem = mapData[i];
        let logo = "";
        if(mapItem.vendor == "Craigslist"){
            
            logo = "/static/img/craigslist_logo.png";
        }
        
        let myrep = await fetch('/static/pages/popup.html');
        let popup_text = "";
        if(myrep.ok){
            popup_text = await myrep.text();
        }
        
       let domparser = new DOMParser();
       let elem = domparser.parseFromString(popup_text,'text/html');
        elem.getElementById("Vendor").src = logo;
        elem.getElementById("Title").innerHTML = mapItem.title;
        elem.getElementById("Price").innerHTML = mapItem.price;
        elem.getElementById("Posted").innerHTML = mapItem.posted;
        
        
        let loc = L.latLng(mapItem.latitude, mapItem.longitude);
        let listing = L.marker()
        .bindPopup(elem.getElementById("body"))
        .setLatLng(loc).addTo(markers);
    }
    markers.addTo(map);
    document.getElementById("spin").style.display = "none";

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
