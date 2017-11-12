function removeImage(event) {
    var imgElem = event.currentTarget.parentNode.parentNode.parentNode; // 지워야하는 element의 부모
    console.log(imgElem.id)
    
    $.ajax({
        url:'removeImage',
        dataType:'json',
        type:'POST',
        data:{'data': imgElem.id},
        success:function(result){
            if(result['result']==true){
                imgElem.parentNode.removeChild(imgElem);
            }
        }
    });
}