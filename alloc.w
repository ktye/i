
// (bytes -- bucket-type)
(f1 $bk 
    var r
    32i  con 7  4i .x * + i32.clz - :.r  
    4i < if (result i32)
      4i
    else
      .r
    end )


// (bytes -- addr)
(f1 $mk
    var t i a b c d
    .x call $bk :.t 
    4i * :i 
    
    while .i i. =0
    do    .i 128 >=s if ! end .i 4i + :i
    
    
    local.set 2
    i32.const 4
    i32.const 0
    i32.load
    i32.mul
    local.set 3
    block  ;; label = @1
      loop  ;; label = @2
        local.get 2
        i32.load
        i32.eqz
        i32.eqz
        br_if 1 (;@1;)
        local.get 2
        local.get 3
        i32.ge_s
        if  ;; label = @3
          unreachable
        end
        local.get 2
        i32.const 4
        i32.add
        local.set 2
        br 0 (;@2;)
      end
    end
    local.get 2
    i32.load
    local.set 4
    local.get 2
    local.get 4
    i32.load
    i32.store
    local.get 2
    i32.const 4
    i32.sub
    local.set 5
    block  ;; label = @1
      loop  ;; label = @2
        local.get 5
        i32.const 4
        local.get 1
        i32.mul
        i32.ge_u
        i32.eqz
        br_if 1 (;@1;)
        local.get 4
        i32.const 1
        local.get 5
        i32.const 2
        i32.shr_u
        i32.shl
        i32.add
        local.set 6
        local.get 6
        local.get 5
        i32.load
        i32.store
        local.get 5
        local.get 6
        i32.store
        local.get 5
        i32.const 4
        i32.sub
        local.set 5
        br 0 (;@2;)
      end
    end
    local.get 4
    i32.const 1
    i32.store
    local.get 4
    i32.const 4
    i32.add
    local.get 0
    i32.store
    local.get 4)

