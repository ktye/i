# patch out built-in math library, replace with -lm */

BEGIN                {x=1;    }
/double hypot_\(/    {x=match($0,/;/);next} #cut both declaration&function
/void cosin_\(/      {x=match($0,/;/);next}
/void sin1\(/        {x=match($0,/;/);next}
/void cos1\(/        {x=match($0,/;/);next}
/double atan2_\(/    {x=match($0,/;/);next}
/double atan_\(/     {x=match($0,/;/);next}
/double satan\(/     {x=match($0,/;/);next}
/double xatan\(/     {x=match($0,/;/);next}
/double exp_\(/      {x=match($0,/;/);next}
/double expmulti\(/  {x=match($0,/;/);next}
/double ldexp_\(/    {x=match($0,/;/);next}
/int32_t frexp1\(/   {x=match($0,/;/);next}
/double frexp2\(/    {x=match($0,/;/);next}
/int64_t frexp3\(/   {x=match($0,/;/);next}
/double normalize\(/ {x=match($0,/;/);next}
/double log_\(/      {x=match($0,/;/);next}
/double modabsfi\(/  {x=match($0,/;/);next}
/double pow_\(/      {x=match($0,/;/);next}
/uint64_t se\(/      {x=match($0,/;/);next}
/exp_/  { sub(/exp_/  ,"exp"  )}
/log_/  { sub(/log_/  ,"log"  )}
/pow_/  { sub(/pow_/  ,"pow"  )}
/atan2_/{ sub(/atan2_/,"atan2")}
/hypot_/{gsub(/hypot_/,"hypot")}
{        if          (x) print}
/^}/                 {x=1;    }
