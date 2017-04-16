// pull in desired CSS/SASS files
require( './styles/main.scss' );
require( 'font-awesome-webpack');
require( '../../node_modules/toastr/toastr.scss')
var $ = jQuery = require( '../../node_modules/jquery/dist/jquery.js' );
require( '../../node_modules/bootstrap-sass/assets/javascripts/bootstrap.js' );
var toastr = require( '../../node_modules/toastr/toastr.js')

// inject bundled Elm app into div#main
var Elm = require( '../elm/Main' );
var app = Elm.Main.embed( document.getElementById( 'main' ), {
    apiUrl: '/api/v1/'
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

