function showOutput() {
    document.getElementById('output-modal').style.display='block';
}

function hideOutput() {
    document.getElementById('output-modal').style.display='none';
}

function showPlayerSettings() {
    document.getElementById('player-settings-modal').style.display='block';
}

function hidePlayerSettings() {
    document.getElementById('player-settings-modal').style.display='none';
}

function showProfiles() {
    document.getElementById('profiles-modal').style.display='block';
}

function hideProfiles() {
    document.getElementById('profiles-modal').style.display='none';
}

function addOutput(message) {
    let output = document.getElementById("output");

    output.value = output.value + message + '\n';
}

function clearOutput() {
    document.getElementById('output').value = "";
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

function playerSelected() {
    let videoPlayerRadio = document.querySelector('input[name="video-player"]:checked');
    let videoPlayer = null;

    if (videoPlayerRadio) {
        videoPlayer = videoPlayerRadio.value;
    } else {
        videoPlayer = "browser";
    }

    if (videoPlayer === "browser") {
        document.getElementById('browser-player-container').style.display = 'block';
    } else {
        document.getElementById('browser-player-container').style.display = 'none';
    }
}

function getBasicSettingsPayload() {
    let encodingPresetDropdown = document.getElementById("encoding-preset");
    let payload = null;

    try {
        payload = {
            "scaleFactor": parseFloat(document.getElementById("scale-factor").value),
            "encodingPreset": encodingPresetDropdown.options[encodingPresetDropdown.selectedIndex].value,
            "qualityTarget": parseInt(document.getElementById("quality-target").value),
            "saturation": parseFloat(document.getElementById("saturation").value),
            "contrast": parseFloat(document.getElementById("contrast").value),
            "brightness": parseFloat(document.getElementById("brightness").value),
            "gamma": parseFloat(document.getElementById("gamma").value),
            "playAfter": document.getElementById("play-after").checked
        };
    } catch(err) {
        addOutput("Could not create request payload, err: " + err);
    }

    return payload;
}

function streamVideo(videoName) {
    getVideoDetails(videoName, function(video, width, height) {
        //width = parseInt(width)
        //height = parseInt(height)

        //let browserPlayerContainer = document.getElementById('browser-player-container');
        let videoPlayer = document.getElementById('video-player');
        //const videoPlayerParent = videoPlayer.parentElement;
        
        //browserPlayerContainer.style.display = "block";
        //const ratio = width / height;
        //const styles = window.getComputedStyle(videoPlayerParent);
        //const containerWidth = videoPlayerParent.clientWidth - parseFloat(styles.paddingLeft) - parseFloat(styles.paddingRight) - parseFloat(styles.marginLeft) - parseFloat(styles.marginRight);
        //const containerWidth = browserPlayerContainer.clientWidth - parseFloat(styles.paddingLeft) - parseFloat(styles.paddingRight);
        //const containerWidth = browserPlayerContainer.clientWidth;

        //videoPlayer.width = containerWidth - 5;
        //videoPlayer.height = containerWidth / ratio;
        document.getElementById('video-title').innerText = video;
        videoPlayer.src = '{{.FrontendUri}}/streamvideo/' + encodeURIComponent(video);
        document.getElementById("available-videos").value = video;
    });
}

function playVideo(videoName) {
    let dropdown = document.getElementById("available-videos");
    let video = "";

    if (videoName) {
        video = videoName;
    } else {
        if (dropdown.options[dropdown.selectedIndex]) {
            video = dropdown.options[dropdown.selectedIndex].value;
        }
    }

    if (video) {
        let videoPlayerRadio = document.querySelector('input[name="video-player"]:checked');
        let videoPlayer = null;

        if (videoPlayerRadio) {
            videoPlayer = videoPlayerRadio.value;
        } else {
            videoPlayer = "browser";
        }

        if (videoPlayer === "browser") {
            streamVideo(video);
            return;
        }

        let payload = {
            "video": video,
            "alternatePlayer": document.getElementById("alternate-player").value
        };

        fetch('{{.FrontendUri}}/playvideo/', {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(function(response) {
            return response.json();
        })
        .then(function(responseJson) {
            if (responseJson.Error) {
                addOutput("Play error: " + responseJson.Error);
            }
        });
    }
}

function deleteVideo() {
    let dropdown = document.getElementById("available-videos");
    let video = dropdown.options[dropdown.selectedIndex];

    if (video) {
        let confirmed = confirm("Are you sure?");

        if (!confirmed) {
            return;
        }

        let payload = {
            "video": video.value
        };

        disableButtons();
        fetch('{{.FrontendUri}}/deletevideo/', {
            method: "DELETE",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(function(response) {
            enableButtons();
            return response.json();
        })
        .then(function(responseJson) {
            if (responseJson.Error) {
                addOutput("Delete video error: " + responseJson.Error);
                }else {
                addOutput(video.value + " was deleted!");
                getAvailableVideos();
            }
        });
    }
}

function clipVideo() {
    if (document.getElementsByClassName("error").length > 0) {
        addOutput("You have clip settings errors, fix them!");
        return;
    }

    let videoDropdown = document.getElementById("available-videos");
    let selectedVideo = videoDropdown.options[videoDropdown.selectedIndex];

    if (!selectedVideo) {
        addOutput("You don't have a video selected.");
        return;
    }

    let payload = getBasicSettingsPayload();
    payload.video = selectedVideo.value;
    payload.startTime = document.getElementById("start-time").value;
    payload.endTime = document.getElementById("end-time").value;
    payload.alternatePlayer = document.getElementById("alternate-player").value;

    disableButtons();
    fetch('{{.FrontendUri}}/clipvideo', {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    })
    .then(function(response) {
        enableButtons();
        return response.json();
    })
    .then(function(responseJson) {
        if (responseJson.Error) {
            addOutput("Clip error: " + responseJson.Error);
        } else {
            addOutput(responseJson.newVideoName + " was created!");
            getAvailableVideos();
            
            if (payload.playAfter) {
                addOutput("Playing: " + responseJson.newVideoName);
                playVideo(responseJson.newVideoName);
            }
        }
    });
}

function getAvailableVideos() {
    let dropdown = document.getElementById("available-videos");
    let selected = "";
    let selectedOption = dropdown.options[dropdown.selectedIndex];

    if (selectedOption) {
        selected = selectedOption.value;
    }

    disableButtons();
    fetch('{{.FrontendUri}}/getavailablevideos')
    .then(function(response) {
        enableButtons();
        return response.json();
    })
    .then(function(responseJson) {
        if (responseJson.Error) {
            addOutput("Get videos error: " + responseJson.Error);
        } else {
            dropdown.innerHTML = "";

            for (let i = 0; i < responseJson.AvailableVideos.length; i++) {
                let videoName = responseJson.AvailableVideos[i];
                let option = document.createElement("option");
                option.text = videoName;
                option.value = videoName;

                if (selected === videoName) {
                    option.selected = "selected";
                }

                dropdown.add(option);
            }

            dropdown.setAttribute('size', responseJson.AvailableVideos.length);
            addOutput("Available videos retrieved");
        }
    });
}

function getVideoDetails(videoName, callback) {
    let dropdown = document.getElementById("available-videos");
    let video = "";

    if (videoName) {
        video = videoName;
    } else {
        if (dropdown.options[dropdown.selectedIndex]) {
            video = dropdown.options[dropdown.selectedIndex].value;
        }
    }

    if (video) {
        disableButtons();
        fetch('{{.FrontendUri}}/getvideodetails/' + encodeURIComponent(video))
        .then(function(response) {
            enableButtons();
            return response.json();
        })
        .then(function(responseJson) {
            if (responseJson.Error) {
                addOutput("Get video details error: " + responseJson.Error);
            } else {
                let resolution = responseJson.Resolution;

                if (resolution && resolution.includes("x")) {
                    let resolutionFields = resolution.split("x");

                    if (Array.isArray(resolutionFields) && resolutionFields.length > 1) {
                        document.getElementById("source-width").value = resolutionFields[0];
                        document.getElementById("source-height").value = resolutionFields[1];
                        calculateScale();

                        if (callback) {
                            callback(video, resolutionFields[0], resolutionFields[1]);
                        }

                        addOutput("Video details retrieved for " + video);
                    }
                }
            }
        });
    }
}

function getConfig() {
    let dropdown = document.getElementById("profiles");
    let selected = "";
    let selectedOption = dropdown.options[dropdown.selectedIndex];

    if (selectedOption) {
        selected = selectedOption.value;
    }

    disableButtons();
    fetch('{{.FrontendUri}}/getconfig')
    .then(function(response) {
        enableButtons();
        return response.json();
    })
    .then(function(responseJson) {
        if (responseJson.Error) {
            addOutput("Get config error: " + responseJson.Error);
        } else {
            let dropdown = document.getElementById("profiles");
            dropdown.innerHTML = "";

            for (let i = 0; i < responseJson.profiles.length; i++) {
                let option = document.createElement("option");
                option.text = responseJson.profiles[i].profileName;
                option.value = responseJson.profiles[i].profileName;

                if (selected === responseJson.profiles[i].profileName) {
                    option.selected = "selected";
                }

                dropdown.add(option);
            }

            dropdown.setAttribute('size', responseJson.profiles.length);
            profiles = responseJson.profiles;
            profileSelected();
            addOutput("Config retrieved");
        }
    });
}

function saveProfile() {
    let profile = getSelectedProfile();

    if (profile) {
        let confirmed = confirm("Are you sure?");

        if (confirmed) {
            doSaveProfile(profile.profileName);
        }
    }
}

function saveNewProfile() {
    let profileName = prompt("Profile Name?");

    if (profileName) {
        let confirmed = confirm("Are you sure?");

        if (confirmed) {
            doSaveProfile(profileName);
        }
    }
}

function deleteProfile() {
    let profile = getSelectedProfile();

    if (profile) {
        let confirmed = confirm("Are you sure?");

        if (confirmed) {
            let payload = {"profileName": profile.profileName};
            
            disableButtons();
            fetch('{{.FrontendUri}}/deleteprofile', {
                method: "DELETE",
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            })
            .then(function(response) {
                enableButtons();
                return response.json();
            })
            .then(function(responseJson) {
                if (responseJson.Error) {
                    addOutput("Delete profile error: " + responseJson.Error);
                } else {
                    addOutput(profile.profileName + " was deleted!");
                    getConfig();
                }
            });
        }
    }
}

function doSaveProfile(profileName) {
    if (profileName) {
        let payload = getBasicSettingsPayload();
        payload.profileName = profileName;
        payload.alternatePlayer = document.getElementById("alternate-player").value;
        payload.videoPlayer = document.querySelector('input[name="video-player"]:checked').value;

        disableButtons();
        fetch('{{.FrontendUri}}/saveprofile', {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(function(response) {
            enableButtons();
            return response.json();
        })
        .then(function(responseJson) {
            if (responseJson.Error) {
                addOutput("Save profile error: " + responseJson.Error);
            } else {
                addOutput(profileName + " was saved!");
                getConfig();
            }
        });
    }
}

function getSelectedProfile() {
    let dropdown = document.getElementById("profiles");
    let selectedOption = dropdown.options[dropdown.selectedIndex];
    let profile = null;

    if (selectedOption) {
        let selected = dropdown.options[dropdown.selectedIndex].value;

        for (let i = 0; i < profiles.length; i++) {
            if (selected === profiles[i].profileName) {
                profile = profiles[i];
                break;
            }
        }
    }

    return profile;
}

function profileSelected() {
    let profile = getSelectedProfile();

    if (profile) {
        document.getElementById("scale-factor").value = profile.scaleFactor;
        document.getElementById("encoding-preset").value = profile.encodingPreset;
        document.getElementById("quality-target").value = profile.qualityTarget;
        document.getElementById("saturation").value = profile.saturation;
        document.getElementById("contrast").value = profile.contrast;
        document.getElementById("brightness").value = profile.brightness;
        document.getElementById("gamma").value = profile.gamma;
        document.getElementById("play-after").checked = profile.playAfter;

        if (profile.videoPlayer === "browser") {
            document.getElementById("player-browser").checked = true;
        } else if (profile.videoPlayer === "ffplay") {
            document.getElementById("player-ffplay").checked = true;
        } else if (profile.videoPlayer === "alternate") {
            document.getElementById("player-alternate").checked = true;
        }

        document.getElementById("alternate-player").value = profile.alternatePlayer;
        playerSelected();
    }
}

function validateTime(event) {
    let value = event.target.value;

    if (typeof value !== "string") {
        event.target.classList.add("error");
        return;
    }

    let values = value.split(":");

    if (values.length !== 3) {
        event.target.classList.add("error");
        return;
    }

    for (let i = 0; i < values.length; i++) {
        if (values[i].length !== 2) {
            event.target.classList.add("error");
            return;
        }

        if (values[i].includes('-')) {
            event.target.classList.add("error");
            return;
        }

        if (isNaN(values[i])) {
            event.target.classList.add("error");
            return;
        }
    }

    event.target.classList.remove("error");
    return;
}

function validateFloat(event, min, max) {
    let value = event.target.value;

    if (isNaN(value)) {
        event.target.classList.add("error");
        return;
    }

    if (value < min || value > max) {
        event.target.classList.add("error");
        return;
    }

    event.target.classList.remove("error");
    return;
}

function validateInt(event, min, max) {
    let value = event.target.value;

    if (isNaN(value) || value.includes('.')) {
        event.target.classList.add("error");
        return;
    }

    value = parseInt(value);

    if (!Number.isInteger(value)) {
        event.target.classList.add("error");
        return;
    }

    if (value < min || value > max) {
        event.target.classList.add("error");
        return;
    }

    event.target.classList.remove("error");
    return;
}

function calculateScale() {
    let scaleFactor = document.getElementById("scale-factor").value;
    let sourceWidth = document.getElementById("source-width").value;
    let sourceHeight = document.getElementById("source-height").value;
    let newWidth = sourceWidth / scaleFactor;
    let newHeight = sourceHeight / scaleFactor;
    let newSize = Math.trunc(newWidth) + "x" + Math.trunc(newHeight);

    document.getElementById("new-size").value = newSize;
}

function disableButtons() {
    let buttons = document.getElementsByTagName("button");
    for (let i = 0; i < buttons.length; i++) {
        buttons[i].disabled = true;
    }
}

function enableButtons() {
    let buttons = document.getElementsByTagName("button");
    for (let i = 0; i < buttons.length; i++) {
        buttons[i].disabled = false;
    }
}

let profiles;

disableButtons();
fetch('{{.FrontendUri}}/checkffmpeg')
    .then(function(response) {
        enableButtons();
        return response.json();
    })
    .then(function(responseJson) {
        if (responseJson.Error) {
                addOutput("Check FFmpeg error: " + responseJson.Error);
        } else {
            let message = "";

            if (responseJson.FFmpegExists === false ||
                responseJson.FFplayExists === false ||
                responseJson.FFprobeExists === false
            ) {
                message = 'Either FFmpeg or FFplay or FFprobe is missing, please download official releases and include them in the same directory as ffmpeg-clipper or add them to your OS PATH. https://ffmpeg.org/download.html';
            }

            if (message !== "") {
                addOutput(message);
            } else {
                getConfig();
                getAvailableVideos();
            }
        }
    }
);