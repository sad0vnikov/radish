// pull in desired CSS/SASS files
require( './styles/main.scss' );
require( 'font-awesome-webpack');
require( '../../node_modules/toastr/toastr.scss')
var $ = jQuery = require( '../../node_modules/jquery/dist/jquery.js' );
require( '../../node_modules/bootstrap-sass/assets/javascripts/bootstrap.js' );
var toastr = require( '../../node_modules/toastr/toastr.js')

// inject bundled Elm app into div#main
var Elm = require( '../elm/Main' );

var currentURI = window.location.href
currentURI = currentURI.replace("index.html", "")
currentURI = currentURI.replace("#", "")
currentURI = currentURI.replace("?", "")

if (currentURI.charAt(currentURI.length - 1) != "/") {
    currentURI += "/";
} 

var xhr = new XMLHttpRequest();
xhr.open("GET", currentURI + 'api/v1/appVersion', false);
xhr.send();

appVersion = "UNKNOWN";
if (xhr.status == 200) {
    console.log(xhr.responseText);
    versionInfo = JSON.parse(xhr.responseText);
    if (versionInfo && versionInfo.Version) {
        appVersion = versionInfo.Version
    }
}

var app = Elm.Main.embed( document.getElementById( 'main' ), {
    apiUrl: currentURI + 'api/v1',
    appVersion: appVersion,
} );

toastr.options.positionClass = "toast-bottom-right";

app.ports.toastError.subscribe(function(message) {
    toastr.error(message)
})

app.ports.toastInfo.subscribe(function(message) {
    toastr.info(message)
})

app.ports.toastWarning.subscribe(function(message) {
    toastr.warning(message)
})

app.ports.toastSuccess.subscribe(function(message) {
    toastr.success(message)
})

app.ports.showConfirmationDialog.subscribe(function(message) {
    $('#modal .modal-body').html(message);
    $('#modal').modal();

    $('#modal .confirm').click(function() {
        $('#modal').data("message", "ok");
        $('#modal').modal('hide');
    });

    $('#modal .cancel').click(function() {
        $('#modal').data("message", "cancel");
    });

    $('#modal').on('hidden.bs.modal', function() {
        app.ports.dialogClosed.send($('#modal').data('message'));
    })
});