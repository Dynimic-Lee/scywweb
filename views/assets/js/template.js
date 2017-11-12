function InitTemplate(uri) {
    console.log("InitTemplate ["+uri+"]");
    switch(uri) {
        case "/" :
        case "/show" : {
            $('#nav-show').addClass("active");
            break;
        }
        case "/upload" : {
            $('#nav-upload').addClass("active");
            break;
        }
        default : {
            break;
        }
    }
}