window.onload = () => {
    setLocale();

    //Intercepts the submit request and does my own post request, doesn't even need a reload.
    document.getElementById("requestForm").addEventListener("submit", async (e) => {
        //set the loading state
        setLoadingstate(true);
        setLocalizedWaitText();

        //get the data from the form
        e.preventDefault();
        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData);

        //disable the form so that the user doesn't do too many requests at once
        disableForm();

        //send the request to the API endpoint
        const response = await fetch("/api/v1/search", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        try {
            //handle the response
            let result = await response.json();

            if (result["result"] === "success") {
                replaceAnswerBlockText(result["content"]);
            } else {
                replaceErrorBlockText(result["content"]);
            }
        } finally {
            //re-enable the form and disable the loading animation
            enableForm();
            setLoadingstate(false);
        }
    });

    document.getElementById("update-button").addEventListener("click", async () => {
        await fetch("/api/v1/update")
            .then(response => {
                if (response.ok) {
                    alert("ok");
                } else {
                    alert("error");
                }
            });
    });

    document.getElementById("reset-button").addEventListener("click", async () => {
        await fetch("/api/v1/reset")
            .then(response => {
                if (response.ok) {
                    alert("ok");
                } else {
                    alert("error");
                }
            });
    });

    document.getElementById("reconfig-button").addEventListener("click", async () => {
        await fetch("/api/v1/reconfig")
            .then(response => {
                if (response.ok) {
                    alert("ok");
                } else {
                    alert("error");
                }
            });
    });
}

//change the text from the answer block, handles markdown
function replaceAnswerBlockText(new_text) {
    document.getElementById("bot-answer").innerHTML = marked.parse(new_text);
}

//show the error message in the error block
function replaceErrorBlockText(new_text) {
    document.getElementById("request-error").innerHTML = new_text;
}

//enable the form for modification
function enableForm() {
    let form = document.getElementById("requestForm")
    let elements = form.elements
    for (let i = 0; i < elements.length; i++) {
        elements[i].disabled = false
    }
}

//disable the form
function disableForm() {
    let form = document.getElementById("requestForm")
    let elements = form.elements
    for (let i = 0; i < elements.length; i++) {
        elements[i].disabled = true
    }
}

function setLoadingstate(is_loading) {
    let bar = document.getElementById("loading-bar");
    if (is_loading) {
        bar.classList.add("animate");
    } else {
        bar.classList.remove("animate");
    }
}
