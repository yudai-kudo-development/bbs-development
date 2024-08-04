function toggleEditForm(topicId) {
    var form = document.getElementById('edit-form-' + topicId);
    if (form.style.display === 'none') {
        form.style.display = 'block';
    } else {
        form.style.display = 'none';
    }
}

function confirmDeletion(topicId) {
    if (confirm('このトピックを削除してもよろしいですか？')) {
        document.getElementById('delete-form-' + topicId).submit();
    }
}
