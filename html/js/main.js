function getSelectedVideo() {
    let videoDropdown = document.getElementById("available-videos-select");
    return videoDropdown.options[videoDropdown.selectedIndex];
}

function getSelectedProfile() {
    let profileDropdown = document.getElementById("profiles-select");
    return profileDropdown.options[profileDropdown.selectedIndex];
}