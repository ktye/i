BEGIN{x=-1}                     #skip header
/^static void panic/{x=1;next}  #last header line

#remove these functions:
/^static int32_t bucket\(.*{$/{x=0}

/^static int32_t alloc\(.*{$/{x=0}
/^static void mfree\(.*{$/{x=0}

{if(x>0)print}
/^}/{if(!x)x=1}                 #end of function
