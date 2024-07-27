function getSelectedVideo() {
    let videoDropdown = document.getElementById("available-videos-select");
    return videoDropdown.options[videoDropdown.selectedIndex];
}

function getSelectedProfile() {
    let profileDropdown = document.getElementById("profiles-select");
    return profileDropdown.options[profileDropdown.selectedIndex];
}

function formatFFmpegTime(time) {
    time = time.toString();

    if (time.length < 2) {
        time = '0' + time;
    }

    return time;
}

function getCurrentFFmpegTime() {
    let videoPlayer = document.getElementById('video-player');
    let currentTime = videoPlayer.currentTime;
    let ffmpegTime = null;
    let hours = 0;
    let minutes = 0;
    let seconds = 0;

    if (currentTime > 0) {
        if (currentTime < 60) {
            seconds = parseInt(currentTime);
        } else if (currentTime < 3600) {
            minutes = parseInt(currentTime / 60);
            seconds = parseInt(currentTime % 60);
        } else {
            hours = parseInt(currentTime / 3600);
            let remainingSeconds = parseInt(currentTime % 3600);

            if (remainingSeconds < 60) {
                seconds = remainingSeconds;
            } else {
                minutes = parseInt(remainingSeconds / 60);
                seconds = parseInt(remainingSeconds % 60);
            }
        }
    }

    ffmpegTime = formatFFmpegTime(hours) + ':' + formatFFmpegTime(minutes) + ':' + formatFFmpegTime(seconds);

    return ffmpegTime;
}

function setClipStart() {
    let ffmpegTime = getCurrentFFmpegTime();

    if (ffmpegTime !== null) {
        document.getElementById("start-time").value = ffmpegTime;
    }
}

function setClipStop() {
    let ffmpegTime = getCurrentFFmpegTime();

    if (ffmpegTime !== null) {
        document.getElementById("end-time").value = ffmpegTime;
    }
}

function closeVideo() {
    document.getElementById('video-title').innerHTML = '';
    let videoPlayer = document.getElementById('video-player');
    videoPlayer.pause();
    videoPlayer.src = '';
    videoPlayer.load();
}

function calculateScale() {
    let scaleFactorElement = document.getElementById("scale-factor");

    if (scaleFactorElement) {
        let scaleFactor = scaleFactorElement.value;
        let sourceWidth = document.getElementById("source-width").value;
        let sourceHeight = document.getElementById("source-height").value;
        let newWidth = sourceWidth / scaleFactor;
        let newHeight = sourceHeight / scaleFactor;
        let newSize = Math.trunc(newWidth) + "x" + Math.trunc(newHeight);
    
        document.getElementById("new-size").value = newSize;
    }
}