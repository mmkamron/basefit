var trashs = document.getElementsByClassName("bi-trash");
var edits = document.getElementsByClassName("bi-pen");

for (let trash of trashs) {
    trash.onclick = () => Delete(trash.id);
};

for (let edit of edits) {
    let weight = edit.parentElement.previousElementSibling
    let activity = weight.previousElementSibling
    edit.onclick = () => Input(edit, weight, activity);
};

function Delete(id) {
    fetch('/gym/' + id, {
        method: 'DELETE'
    }).then(() => {
        window.location.reload();
    })
};

function Input(edit, weight, activity) {
    let id = edit.id
    let inputActivity = '<input type="text" name="activity">'
    let inputWeight = '<input type="number" name="weight">'
    let inputEdit = '<button class="btn btn-primary">Submit</button>'
    weight.innerHTML=inputWeight
    activity.innerHTML=inputActivity
    edit.outerHTML=inputEdit
    inputEdit.onclick = () => Update(id, weight.value, activity.value);
} 

function Update(id, activity, weight) {
    console.log(id, activity, weight)
    // fetch('/gym/' + id, {
    //     method: 'PUT'
    // }).then(() => {
    //     window.location.reload();
    // })
};
