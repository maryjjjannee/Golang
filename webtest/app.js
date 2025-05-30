function appendData(data) {
    var mainContainer = document.getElementById('myData');
    for (var i = 0; i < data.length; i++) {
        var div = document.createElement('div');
        div.innerHTML = 'CourseID: ' + data[i].id + ' ' + data[i].name + ' ' + data[i].price + ' ' + data[i].instructor;
        mainContainer.appendChild(div);
    }
}

fetch('http://localhost:8080/course')
  .then(res => {
    if (!res.ok) throw new Error('Network response was not ok');
    return res.json();
  })
  .then(data => {
    appendData(data);
  })
  .catch(err => {
    console.error('Fetch error:', err);
  });