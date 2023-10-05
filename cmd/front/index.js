// const order = document.getElementById('order_uid')
// console.log(order)
const btn = document.getElementById('get_btn')
// const currentUrl = window.location.href;
const serverUrl = "order_uid/"
// document.getElementById('order_uid').


btn.addEventListener('click', getData)
async function getData(e) {
    e.preventDefault()
    const order = document.getElementById('order_uid').value
    console.log(order)
    const resp = await fetch(serverUrl+order)
    console.log(resp)

    // fetch (serverUrl+order, {method: 'GET'})
    // .then(response => response.json())
    // .then(data => console.log(data))
}