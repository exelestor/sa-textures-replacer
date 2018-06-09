document.getElementById('rplcBtn').onclick = function () {
    controller.replace(
        document.getElementById('txdDir').value,
        document.getElementById('picUrl').value
    );
};

document.getElementById('txdDir').onclick = function () {
    external.invoke('opendir')
};

function progressbar(value) {
    let elem = document.getElementById("myBar");
    elem.style.width = value + '%';
    // elem.innerHTML = value * 1 + '%';
}
