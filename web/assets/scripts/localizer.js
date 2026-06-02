
const lang = navigator.language.split("-")[0];

async function setLocale() {
    let queryTextBox = document.getElementById("query");
    let submitButton = document.getElementById("submit-button");
    let answerTextBlock = document.getElementById("bot-answer");

    let locale = await loadLocale(lang);

    queryTextBox.placeholder = locale["placeholder"];
    submitButton.value = locale["submit"];
    answerTextBlock.innerHTML = `<span>${locale["welcome"]}</span>`
}

async function setLocalizedWaitText() {
    let answerTextBlock = document.getElementById("bot-answer");

    let locale = await loadLocale(lang);

    answerTextBlock.innerHTML = `<span>${locale["waiting"]}</span>`
}

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
