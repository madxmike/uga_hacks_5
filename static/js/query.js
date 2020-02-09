let markers = L.layerGroup();
let latlng;

let map = L.map('mapid', {
    minZoom: 8,
    maxZoom: 15
}).setView([33.7490, -84.3880], 11);

L.tileLayer('http://{s}.google.com/vt/lyrs=m&x={x}&y={y}&z={z}', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    subdomains:['mt0','mt1','mt2','mt3']
}).addTo(map);

let platform = new H.service.Platform({
    'apikey': 'eaQF8vvay1leHY9YE66PJkFL3Fh4OwwECClcpFV760Y'
  });

map.on('popupopen', function(e) {
    console.log("open")
    var elements = document.querySelectorAll(".leaflet-tile-container,.leaflet-marker-pane", );
    var blured = Array.prototype.filter.call(elements, function(element){
        return element.style.filter = "blur(2px)"
    });
});

map.on('popupclose', function(e) {
    console.log("close")
    var elements = document.querySelectorAll(".leaflet-tile-container,.leaflet-marker-pane", );
    var blured = Array.prototype.filter.call(elements, function(element){
        return element.style.filter = "blur(0px)"
    });
});

document.getElementById("map").onclick = function () {
    document.getElementById("mapid").style.visibility = "visible";
    document.getElementById("mapid").classList.remove("w-0");
    document.getElementById("list").style.visibility = "hidden";
    
};
document.getElementById("grid").onclick = function () {
    document.getElementById("list").style.visibility = "visible";
    document.getElementById("mapid").classList.add("w-0");
    document.getElementById("mapid").style.visibility = "hidden";
};

search.onsubmit = async (e) => {
    e.preventDefault();
    latlng = map.getCenter();
    document.getElementById("spin").style.display = "inline-block";
    
    let form = new FormData(search);
    let address = form.get("address") + " " + form.get('city') + " " + form.get('region');

    let promise = toLatLng(address)

    await promise.then((message) => {
        setLoc(message)
    }).catch((error) => {
        latlng = map.getCenter();
        console.log(error);
    });
    setDefaults(form);
    
    let maxBounds = latlng.toBounds(form.get("maxDistance")*1609.34);
    form.append("bounds", maxBounds.toBBoxString());
    
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
    let myrep = await fetch('/static/pages/popup.html');
    let popup_text = "";
    if(myrep.ok){
        popup_text = await myrep.text();
    }  
   
    populateMap(mapData, maxBounds, popup_text, markers);
    map.setView(latlng, 9);
    document.getElementById("spin").style.display = "none";
    
  };//onSubmit
  
  function setLoc(result, address) {
      
        let lat = result.response.view[0].result[0].location.displayPosition.latitude;
        let lng = result.response.view[0].result[0].location.displayPosition.longitude;
        latlng = L.latLng(lat, lng);
      
        
      
    
    
    
    
  }//setloc
  
  function setDefaults(form){
    if (form.get("price_min") === "") {
        form.set("price_min", "0")
    }
    if (form.get("price_max") === "") {
        form.set("price_max", "1000")
    }
    if(form.get("maxDistance") == ""){
        form.set("maxDistance", 40);
    }
  }//setDefaults

  function toLatLng(address) {
    
    return new Promise((resolve, reject) => {
        let geocoder = platform.getGeocodingService(),
        geocodingParameters = {
            searchText: address,
            jsonattributes : 1
        };
    geocoder.geocode(
        geocodingParameters,
        resolve,
        reject
      );
    })
 
  }//toLatLang

  function populateMap(mapData, maxBounds, popup_text, markers){
    for(let i=0; i<mapData.length; i++){
       
        let mapItem = mapData[i];
        let listLoc = L.latLng(mapItem.latitude, mapItem.longitude);
        if(!maxBounds.contains(listLoc)){
            continue;
        }
        let logo = "";
        if(mapItem.vendor == "Craigslist"){
            
            logo = "/static/img/craigslist_logo.png";
        }
        let domparser = new DOMParser();
        let elem = domparser.parseFromString(popup_text,'text/html');
        elem.getElementById("Vendor").src = logo;
        elem.getElementById("Title").innerHTML = mapItem.title;
        elem.getElementById("Price").innerHTML = "$" + mapItem.price;
        elem.getElementById("Posted").innerHTML = mapItem.posted;
        elem.getElementById("Link").href = mapItem.url;
        elem.getElementById("Description").innerHTML = mapItem.description;
        
        let loc = L.latLng(mapItem.latitude, mapItem.longitude);
        let listing = L.marker()
        .bindPopup(elem.getElementById("body"), {maxWidth: 700, maxHeight: 500})
        .bindTooltip(mapItem.title)
        .setLatLng(loc)
        .addTo(markers);
        
        document.getElementById("list").append(elem.getElementById("body"));
        
     }
     markers.addTo(map);
    }
    
  function onError(error) {
    alert('Can\'t load address.');

  }
