
//Get the language locale from the navigator
const lang = navigator.language.split("-")[0];

//Set the localized texts for the index.
//I do it in the Js because it's not that much text we have to serve
//And also we have a changing text, which would be annoying to manage from the server.
async function setLocale() {
    let queryTextBox = document.getElementById("query");
    let submitButton = document.getElementById("submit-button");
    let answerTextBlock = document.getElementById("bot-answer");

    let locale = await loadLocale(lang);

    queryTextBox.placeholder = locale["placeholder"];
    submitButton.value = locale["submit"];
    //Set a span because this is the block for the Markdown
    //So it doesn't have any text inside
    answerTextBlock.innerHTML = `<span>${locale["welcome"]}</span>`
}

//Function to set the localized waiting text
async function setLocalizedWaitText() {
    let answerTextBlock = document.getElementById("bot-answer");

    let locale = await loadLocale(lang);

    answerTextBlock.innerHTML = `<span>${locale["waiting"]}</span>`
}

//Function to load the locale file containing the actual text.
async function loadLocale(lang) {
    try {
        const res = await fetch(`/assets/locales/index/${lang}.json`);
        if (!res.ok) throw new Error();
        return await res.json();
    } catch {
        const res = await fetch(`/assets/locales/index/en.json`);
        return await res.json()
    }
}
