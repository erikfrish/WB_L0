const btn = document.getElementById('get_btn')
const serverUrl = "order_uid/"

btn.addEventListener('click', getData)
async function getData(e) {
    e.preventDefault()
    const order = document.getElementById('order_uid').value
    const resp_form = document.getElementById('response')
    console.log(order)

    fetch(serverUrl+order).then(function (response) {
        response.json().then(function (json) {
            res = JSON.stringify(json, null, 2) 
            resp_form.innerHTML = res;
            console.log(res.getData());
        });
      });
}