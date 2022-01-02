var pathMappings = []
var keyMappings = []

var previousSelectedRow = undefined

function showPane(paneIdQuery){

    document.querySelector(paneIdQuery).classList.add("is-visible")

}

function hideSelf(event){

    event.target.classList.remove("is-visible")

}

function toggleSelf(event) {


    function getNumberOfSelectedItems(parentNode) {


        let siblings = parentNode.children
        let numberOfSelectedItems = 0

        for (let i = 0; i < siblings.length; i++) {

            let isVisible = siblings.item(i).classList.contains("is-selected")

            if (isVisible) {

                numberOfSelectedItems++

            }
        }

        return numberOfSelectedItems

    }

    function getSelectedItem(parentNode) {

        let siblings = parentNode.children

        for (let i = 0; i < siblings.length; i++) {

            let node = siblings.item(i)

            if (node.classList.contains("is-selected")){

                return node

            }

        }


        return null


    }



    if (getNumberOfSelectedItems(event.target.parentNode) === 1){

        let selectedItem = getSelectedItem(event.target.parentNode)

        selectedItem.classList.remove("is-selected")

        setMenuVisibility(selectedItem, false)

        event.target.classList.add("is-selected")

        setMenuVisibility(event.target, true)

    }

    // Reset delete button
    document.querySelectorAll(".del").forEach(e => e.classList.remove("is-active"))


}

function setMenuVisibility (itemNode, visible) {


    let id = itemNode.attributes.getNamedItem("id").value

    let menu = document.querySelector(`div[id="${id.replace('item', 'menu')}"]`)

    if (visible) {

        menu.classList.add("is-visible")
    }

    else {
        menu.classList.remove("is-visible")
    }

}

function toggleRow(event) {

    let children = event.target.parentElement.parentElement.children

    // Deselect all elements
    for (let i = 0; i < children.length; i++) {

        let row = event.target.parentElement.parentElement.children[i]

        for (let j = 0; j < row.children.length ; j++) {

            let cell = row.children[j]

            cell.classList.remove("is-selected")

        }

    }

    // Select element

    for (let i = 0; i < event.target.parentElement.children.length; i++) {

        event.target.parentElement.children[i].classList.add("is-selected")

    }


    event.target.parentElement.parentElement.parentElement.querySelector(".del").classList.add("is-active")

}

function search(event) {

    let table = event.target.parentElement.parentElement.previousElementSibling
    let searchField = document.getElementById(event.target.id)
    let value = searchField.value

    let pattern = value !== undefined ? ".*" + value + ".*" : ".*"

    let regex = new RegExp(pattern)

    // Search in path mapping table
    if (/^pm.*/.test(event.target.id)) {

        let matchedPathMappings = []

        for (let i = 0; i < pathMappings.length; i++) {

            let objValue = JSON.stringify(pathMappings[i])
                .replaceAll("{", "")
                .replaceAll("}", "")
                .replaceAll("\"", "")
                .replaceAll(":", "")
                .replaceAll(",", "")
                .replaceAll(" ", "")

            if (regex.test(objValue)) {

                matchedPathMappings.push(pathMappings[i])

            }

        }

        populateTableWithData(matchedPathMappings, table, "path-mapping")

    }

    // Search in key mapping table
    else if (/^km.*/.test(event.target.id)) {

        let matchedKeyMappings = []

        for (let i = 0; i < keyMappings.length; i++) {

            let objValue = JSON.stringify(keyMappings[i])
                .replaceAll("{", "")
                .replaceAll("}", "")
                .replaceAll("\"", "")
                .replaceAll(":", "")
                .replaceAll(",", "")
                .replaceAll(" ", "")

            if (regex.test(objValue)) {

                matchedKeyMappings.push(keyMappings[i])

            }

        }

        populateTableWithData(matchedKeyMappings, table, "key-mapping")



    }


}

function populateRateLimitMenu() {}

function populatePathMappingsMenu() {

    const url = new Request("http://localhost:8080/config/path-mappings")

    pathMappings = []


    fetch(url)
        .then(response => {return response.json()})
        .then(response => {

            document.querySelector("#pm-table-body").innerHTML = ''

            response.response.forEach(datum => {

                let id = datum['id']
                let path = datum['path']
                let table = datum['table']

                pathMappings.push({"id": id, "path": path, "table": table})

            })

            populateTableWithData(pathMappings, document.querySelector("#pm-table-body"), "path-mapping")


        })



}

function populateTableWithData(data, element, dataType) {

    element.innerHTML = ""

    if (dataType.toLowerCase() === "path-mapping") {

        for (let i = 0; i < data.length; i++) {

            let datum = data[i]

            let id = datum['id']
            let path = datum['path']
            let table = datum['table']



            element.innerHTML += `
    
                    <div class="row">
            
                        <input class="id-cell" readonly type="text" value="${id}" onfocus="toggleRow(event)"/>
                        <input class="path-cell" type="text" value="${path}" onfocus="toggleRow(event)"/>
                        <input class="table-cell" type="text" value="${table}" onfocus="toggleRow(event)"/>
            
                    </div>
                
                
                `

        }

    } else if (dataType.toLowerCase() === "key-mapping") {

        for (let i = 0; i < data.length; i++) {

            let datum = data[i]

            let id = datum['id']
            let key = datum['key']
            let column = datum['column']



            element.innerHTML += `
    
                    <div class="row">
            
                        <input class="id-cell" readonly type="text" value="${id}" onfocus="toggleRow(event)"/>
                        <input class="key-cell" type="text" value="${key}" onfocus="toggleRow(event)"/>
                        <input class="column-cell" type="text" value="${column}" onfocus="toggleRow(event)"/>
            
                    </div>
                
                
                `

        }

    }







}

function populateKeyMappingsMenu() {

    const url = new Request("http://localhost:8080/config/key-mappings")

    keyMappings = []

    fetch(url)
        .then(response => {return response.json()})
        .then(response => {

            document.querySelector("#km-table-body").innerHTML = ''

            response.response.forEach(datum => {

                let id = datum['id']
                let key = datum['key']
                let column = datum['column']

                keyMappings.push({"id": id, "key": key, "column": column})

            })


            populateTableWithData(keyMappings, document.querySelector("#km-table-body"), "key-mapping")

        })



}

function populateBehaviorsMenu() {}

function populateDBConfigMenu() {}

function deleteEntry(event, tableName) {

    let table = event.target.parentElement.parentElement.previousElementSibling

    let selectedRows = table.querySelectorAll(".is-selected")

    if (tableName.toLowerCase() === "path-mappings") {

        let id = undefined
        let path = undefined
        let table = undefined

        for (let i = 0; i < selectedRows.length; i++) {

            for (let j = 0; j < selectedRows[i].classList.length; j++) {

                switch (selectedRows[i].classList[j]) {

                    case "id-cell":

                        id = selectedRows[i].value
                        break

                    case "path-cell":

                        path = selectedRows[i].value
                        break

                    case "table-cell":

                        table = selectedRows[i].value
                        break

                }

            }

        }

        let payload = {

            "Must": {

                "id": id,
                "path": path,
                "table": table

            }

        }

        let url = "http://localhost:8080/config/path-mappings"

        let request = new Request(url, {method: "DELETE", body: JSON.stringify(payload)})

        fetch(request)
            .then(res => {

                if (res.status === 200) {


                    populatePathMappingsMenu()

                }


            })

    }

    else if (tableName.toLowerCase() === "key-mappings") {

        let id = undefined
        let key = undefined
        let column = undefined

        for (let i = 0; i < selectedRows.length; i++) {

            for (let j = 0; j < selectedRows[i].classList.length; j++) {

                switch (selectedRows[i].classList[j]) {

                    case "id-cell":

                        id = selectedRows[i].value
                        break

                    case "key-cell":

                        key = selectedRows[i].value
                        break

                    case "column-cell":

                        column = selectedRows[i].value
                        break

                }

            }

        }

        let payload = {

            "Must": {

                "id": id,
                "key": key,
                "column": column

            }

        }

        let url = "http://localhost:8080/config/key-mappings"

        let request = new Request(url, {method: "DELETE", body: JSON.stringify(payload)})

        fetch(request)
            .then(res => {

                if (res.status === 200) {


                    populateKeyMappingsMenu()

                }


            })


    }






}