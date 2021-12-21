function showAPIPane(){

    document.querySelector("#api-pane").setAttribute("is-visible", "true")

}

function hideAPIPane(event){

    event.target.setAttribute("is-visible", "false")

}