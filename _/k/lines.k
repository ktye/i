/line drawing (k9)
/single data point per column

/ 64*64 to 32*64 rows/cols with half blocks U+2580 U+2584
line32:{o:(32*64)#0        
        o[x+64*2/y]:1+2\y  
        s:(" ";0xe29680;0xe29684)
      ""0:,/(,[;"\n"])',/'s@64^o}


N:64; line32[!N;N/(!N)*!N]


/ 72*72 to 12*72 rows/cols with 6 legacy hlines U+1fb76..U+1fb7b
line12:{o:(12*72)#0        
        o[x+72*6/y]:1+6\y  
        s:(" ";0xf09fadb6;0xf09fadb7;0xf09fadb8;0xf09fadb9;0xf09fadba;0xf09fadbb) 
      ""0:,/(,[;"\n"])',/'s@72^o}


N:72; line12[!N;N/(!N)*!N]


/ 64*64 to 16*64 rows/cols with scan lines U+23bb..U+23bd
line16:{o:(16*64)#0
        o[x+64*4/y]:1+4\y  
        s:" ",3^0xe28ebae28ebbe28ebce28ebd
      ""0:,/(,[;"\n"])',/'s@64^o}


N:64; line16[!N;N/(!N)*!N]


/ 64*64 to 8*64 rows/cols with blocks U+2581..U+U2588 + vfill
line8:{o:( 8*64)#0
       o[x+64*yy:8/y]:1+8\y
       s:" ",3^0xe29688e29687e29686e29685e29684e29683e29682e29681
 o[,/x+64*(1+yy)_\!8]:1
 ""0:,/(,[;"\n"])',/'s@64^o}

w:64; line8[!w;64/(!w)*!w]

\\
