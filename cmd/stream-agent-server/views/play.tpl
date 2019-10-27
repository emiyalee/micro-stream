<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Page Title</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" type="text/css" href="/css/bootstrap-3.3.7-dist/bootstrap.min.css" />
    <script src="/js/jquery/jquery-1.11.1.min.js"></script>
    <script src="/js/bootstrap-3.3.7-dist/bootstrap.min.js"></script>
    <script src="/js/hls-0.9.1-dist/hls.min.js"></script>
</head>

<body>
    <div class="container">
        <h1 class="text-center">Xingtu Streaming System</h1>
        <div class="row">
            <div class="col-lg-3" style="border-right-style: solid">

            </div>
            <div class="col-lg-9" style="border-left-style: solid">
                <video id="video" width="960" height="540" controls src='{{.stream_url}}'></video>
            </div>
        </div>
    </div>

</body>

<script>
    if (Hls.isSupported()) {
        var video = document.getElementById('video');
        var hls = new Hls();
        hls.loadSource(video.src);
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED, function () {
            video.play();
        });
    }
    // hls.js is not supported on platforms that do not have Media Source Extensions (MSE) enabled.
    // When the browser has built-in HLS support (check using `canPlayType`), we can provide an HLS manifest (i.e. .m3u8 URL) directly to the video element throught the `src` property.
    // This is using the built-in support of the plain video element, without using hls.js.
    else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        // video.src = 'http://172.16.5.150:11241/vod/1080p/1080p.m3u8';
        video.addEventListener('canplay', function () {
            video.play();
        });
    }
</script>

</html>