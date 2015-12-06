import d3 from 'd3';


export function formatBytes(bytes) {
    //WRONG
    if(bytes == 0) return '0 Byte';
    var k = 1000;
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)) + ' ' + sizes[i];
}


export function formatPercent(val) {
    return val + "%"; 
}
