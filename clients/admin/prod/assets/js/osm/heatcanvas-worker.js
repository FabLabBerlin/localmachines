function calc(a){value=a.value||{},degree=a.degree||1;for(var b in a.data)for(var c=a.data[b],d=Math.floor(Math.pow(c/a.step,1/degree)),e=Math.floor(b%a.width),f=Math.floor(b/a.width),g=e-d;g<e+d;g+=1)if(!(g<0||g>a.width))for(var h=f-d;h<f+d;h+=1)if(!(h<0||h>a.height)){var i=Math.sqrt(Math.pow(g-e,2)+Math.pow(h-f,2));if(!(i>d)){var j=c-a.step*Math.pow(i,degree),k=g+h*a.width;value[k]?value[k]=value[k]+j:value[k]=j}}postMessage({value:value})}onmessage=function(a){calc(a.data)};