
function emulateKeyPress(keyChar, keyCode) {

    var keyEvent = new KeyboardEvent("keydown", {
        key: keyChar,
        code: keyChar,
        which: keyCode,
        keyCode: keyCode,
        charCode: keyCode,
        bubbles: true,
        cancelable: true,
    });

    var inputElement = document.getElementById("terminal");

    inputElement.dispatchEvent(keyEvent);
}

/*key: 's',
        code: 's',
        which: 83,
        keyCode: 83,
        charCode: 83,
        bubbles: true,
        cancelable: true,*/