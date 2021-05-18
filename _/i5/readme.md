<details><summary><code>+</code>flip add</summary>
<a href="./numeric.go#L280"><code>(1;2) /1 2</code></a><br>
<a href="./numeric.go#L281"><code>(1;2 3) /(1;2 3)</code></a><br>
<a href="./numeric.go#L282"><code>(1)+2 /3</code></a><br>
<a href="./numeric.go#L283"><code>(1+2;4)+3 /6 7</code></a><br>
<a href="./numeric.go#L284"><code>1+1 /2</code></a><br>
<a href="./numeric.go#L285"><code>2 3+6 /8 9</code></a><br>
<a href="./numeric.go#L286"><code>(1 2;3 4)+1 /(2 3;4 5)</code></a><br>
<a href="./numeric.go#L287"><code>1+(2 3;4) /(3 4;5)</code></a><br>
<a href="./numeric.go#L288"><code>(1;2 3)+(2.;4 5) /(3.;6 8)</code></a><br>
<a href="./verbs.go#L175"><code>+(1 2 3;4 5 6) /(1 4;2 5;3 6)</code></a><br>
<a href="./verbs.go#L176"><code>+(1 2;&#34;ab&#34;) /((1;&#34;a&#34;);(2;&#34;b&#34;))</code></a><br>
</details>
<details><summary><code>-</code>neg sub</summary>
<a href="./numeric.go#L294"><code>2 3-6 /-4 -3</code></a><br>
<a href="./numeric.go#L295"><code>(1-2)-3 /-4</code></a><br>
<a href="./numeric.go#L296"><code>1-2-3 /2</code></a><br>
<a href="./numeric.go#L297"><code>1-(2-3) /2</code></a><br>
<a href="./numeric.go#L302"><code>7 8-6 5 /1 3</code></a><br>
</details>
<details><summary><code>*</code>first mul</summary>
<a href="./numeric.go#L308"><code>7*8 /56</code></a><br>
<a href="./verbs.go#L3"><code>*3 2 1 /3</code></a><br>
<a href="./verbs.go#L4"><code>*2 /2</code></a><br>
</details>
<details><summary><code>%</code>sqrt div</summary>
<a href="./numeric.go#L314"><code>%9 /3.</code></a><br>
<a href="./numeric.go#L319"><code>9%2 /4</code></a><br>
<a href="./numeric.go#L320"><code>9%2. /4.5</code></a><br>
</details>
<details><summary><code>&amp;</code>where min</summary>
<a href="./index.go#L106"><code>&amp;2 0 1 /0 0 2</code></a><br>
<a href="./numeric.go#L326"><code>3&amp;4 /3</code></a><br>
<a href="./numeric.go#L327"><code>3 4.&amp;3.5 /3 3.5</code></a><br>
</details>
<details><summary><code>|</code>reverse max</summary>
<a href="./numeric.go#L346"><code>3|4 /4</code></a><br>
<a href="./numeric.go#L347"><code>3 4.|3.5 /3.5 4</code></a><br>
<a href="./verbs.go#L268"><code>|3 2 8 /8 2 3</code></a><br>
</details>
<details><summary><code>!</code>til dict</summary>
<a href="./verbs.go#L62"><code>!5 /0 1 2 3 4</code></a><br>
<a href="./verbs.go#L63"><code>!{a:3;b:4} /`a`b!3 4</code></a><br>
<a href="./verbs.go#L121"><code>2!3 /(,2)!,3</code></a><br>
</details>
<details><summary><code>~</code>not match</summary>
<a href="./find.go#L9"><code>1~2 /0b</code></a><br>
<a href="./find.go#L10"><code>0 1 2~!3 /1b</code></a><br>
<a href="./find.go#L11"><code>(+)~(+) /0b</code></a><br>
<a href="./find.go#L21"><code>~!3 /100b</code></a><br>
<a href="./find.go#L22"><code>~(1;010b) /(0b;101b)</code></a><br>
<a href="./find.go#L441"><code>(/&#34;[1-3]&#34;)~&#34;alpha&#34; /0b</code></a><br>
<a href="./find.go#L442"><code>(/&#34;[1-3]&#34;)~&#34;alpha24&#34; /1b</code></a><br>
</details>
<details><summary><code>,</code>enlist cat</summary>
<a href="./cat.go#L7"><code>,1 /,1</code></a><br>
<a href="./cat.go#L8"><code>,1 2 /,1 2</code></a><br>
<a href="./cat.go#L9"><code>() /()</code></a><br>
<a href="./cat.go#L10"><code>(1;) /(1;)</code></a><br>
<a href="./cat.go#L11"><code>(;1) /(;1)</code></a><br>
<a href="./cat.go#L12"><code>(;;) /(;;)</code></a><br>
<a href="./cat.go#L13"><code>(1;(2;3 4);5) /(1;(2;3 4);5)</code></a><br>
</details>
<details><summary><code>^</code>sort cut</summary>
<a href="./sort.go#L7"><code>^6 3 2 3 /2 3 3 6</code></a><br>
<a href="./sort.go#L8"><code>x:6 3 2 3;^x /2 3 3 6</code></a><br>
<a href="./verbs.go#L203"><code>2^!5 /(0 1;2 3 4)</code></a><br>
<a href="./verbs.go#L204"><code>2 5^&#34;alphabeta&#34; /(&#34;pha&#34;;&#34;beta&#34;)</code></a><br>
<a href="./verbs.go#L205"><code>&#34;abe&#34;^&#34;albetq&#34; /(&#34;al&#34;;,&#34;b&#34;;&#34;etq&#34;)</code></a><br>
</details>
<details><summary><code>=</code>group equal</summary>
<a href="./compare.go#L11"><code>1 2=2 /01b</code></a><br>
<a href="./find.go#L144"><code>=1 3 1 1 5 /1 3 5!(0 2 3;,1;,4)</code></a><br>
</details>
<details><summary><code>&lt;</code>gradeup less</summary>
<a href="./compare.go#L21"><code>2&gt;!4 /1100b</code></a><br>
<a href="./compare.go#L22"><code>`a`b`c&lt;`alpha /100b</code></a><br>
</details>
<details><summary><code>&gt;</code>gradedown more</summary>
<a href="./compare.go#L16"><code>2&lt;!4 /0001b</code></a><br>
</details>
<details><summary><code>#</code>count take</summary>
<a href="./verbs.go#L17"><code>3#1 /1 1 1</code></a><br>
<a href="./verbs.go#L18"><code>&#34;abc&#34;#&#34;ab0dbb&#34; /&#34;abbb&#34;</code></a><br>
<a href="./verbs.go#L50"><code>#222 /1</code></a><br>
<a href="./verbs.go#L51"><code>#!5 /5</code></a><br>
<a href="./verbs.go#L52"><code>#() /0</code></a><br>
</details>
<details><summary><code>_</code>floor drop</summary>
<a href="./numeric.go#L517"><code>_-0.1 3.3 /-1 3</code></a><br>
<a href="./verbs.go#L85"><code>1_2 3 /,3</code></a><br>
<a href="./verbs.go#L86"><code>2_1 2 /!0</code></a><br>
<a href="./verbs.go#L87"><code>3_1 2 /!0</code></a><br>
<a href="./verbs.go#L88"><code>-1_2 3 /,2</code></a><br>
<a href="./verbs.go#L89"><code>-5_2 3 /!0</code></a><br>
<a href="./verbs.go#L90"><code>2 3_!5 /0 1 4</code></a><br>
</details>
<details><summary><code>@</code>typ atx</summary>
<a href="./assign.go#L83"><code>@[1 2 3;1;0] /1 0 3</code></a><br>
<a href="./assign.go#L84"><code>@[1 2 3;1 2;0] /1 0 0</code></a><br>
<a href="./assign.go#L85"><code>@[1 2 3;1 2;4 5] /1 4 5</code></a><br>
<a href="./assign.go#L86"><code>@[1 2 3;1;0 0] /(1;0 0;3)</code></a><br>
<a href="./assign.go#L87"><code>@[1 2 3;1;2.] /(1;2.;3)</code></a><br>
<a href="./assign.go#L88"><code>@[1 2 3;1;+;2] /1 4 3</code></a><br>
<a href="./index.go#L3"><code>4 3 2 1@0 /4</code></a><br>
<a href="./index.go#L4"><code>4 3 2 1@0 3 /4 1</code></a><br>
<a href="./index.go#L5"><code>1 2 3[1] /2</code></a><br>
<a href="./index.go#L6"><code>1 2 4[3-1] /4</code></a><br>
<a href="./index.go#L7"><code>1 2 3  1 /2</code></a><br>
<a href="./index.go#L8"><code>(0.+!10)(1;2 3) /(1.;2 3.)</code></a><br>
<a href="./index.go#L94"><code>`k(&#34;a\nb c&#34;;1 2;3.0 4) /&#34;(\&#34;a\\nb c\&#34;;1 2;3 4.)&#34;</code></a><br>
<a href="./json.go#L24"><code>`json 1 2 3 /&#34;[1,2,3]&#34;</code></a><br>
<a href="./json.go#L25"><code>`json(1;2 3;4.) /&#34;[1,[2,3],4]&#34;</code></a><br>
<a href="./types.go#L404"><code>@1 /`i</code></a><br>
<a href="./types.go#L405"><code>@1 2 /`I</code></a><br>
</details>
<details><summary><code>?</code>uniq fnd</summary>
<a href="./find.go#L77"><code>2 3 4?3 /1</code></a><br>
<a href="./find.go#L78"><code>2 3 4?6 /3</code></a><br>
<a href="./find.go#L79"><code>2 3 4?!5 /3 3 0 1 2</code></a><br>
<a href="./find.go#L80"><code>&#34;a13b145&#34;?(/&#34;[0-9]+&#34;) /(1 2;4 5 6)</code></a><br>
<a href="./find.go#L112"><code>?3 2 3 4 2 /2 3 4</code></a><br>
<a href="./find.go#L113"><code>?(1;&#34;ab&#34;;2;1) /(1;&#34;ab&#34;;2)</code></a><br>
<a href="./json.go#L8"><code>`json?&#34;1 &#34; /1</code></a><br>
</details>
<details><summary><code>$</code>str cast</summary>
<a href="./string.go#L19"><code>$1.23 /&#34;1.23&#34;</code></a><br>
<a href="./string.go#L20"><code>$1 2 3 /(,&#34;1&#34;;,&#34;2&#34;;,&#34;3&#34;)</code></a><br>
<a href="./string.go#L32"><code>`i$&#34;3&#34; /3</code></a><br>
<a href="./string.go#L33"><code>`f$(&#34;3&#34;;&#34;4.5&#34;) /3 4.5</code></a><br>
<a href="./string.go#L34"><code>`z$&#34;3.1&#34; /3.1a0</code></a><br>
<a href="./string.go#L35"><code>`z$&#34;3a20&#34; /3a20</code></a><br>
<a href="./string.go#L36"><code>`$&#34;alpha&#34; /`alpha</code></a><br>
<a href="./string.go#L37"><code>`$(&#34;a&#34;;&#34;bc&#34;) /`a`bc</code></a><br>
<a href="./string.go#L38"><code>0$&#34; ab c  &#34; /&#34;ab c&#34;</code>(trim)</a><br>
<a href="./string.go#L39"><code>3$(&#34;a&#34;;&#34;beta&#34;) /(&#34;a  &#34;;&#34;bet&#34;)</code>(pad)</a><br>
</details>
<details><summary><code>.</code>val call</summary>
<a href="./assign.go#L119"><code>.[(1;2 3);1 0;5] /(1;5 3)</code></a><br>
<a href="./assign.go#L120"><code>.[(1;2 3);1 0;+;5] /(1;7 3)</code></a><br>
<a href="./assign.go#L121"><code>.[(1 2 3;4 5 6);(1;1 2);+;5] /(1 2 3;4 10 11)</code></a><br>
<a href="./assign.go#L122"><code>.[(1 2 3;4 5 6);(0 1;1 2);+;5] /(1 7 8;4 10 11)</code></a><br>
<a href="./assign.go#L123"><code>.[(1 2 3;4 5 6);(;2);+;1] /(1 2 4;4 5 7)</code></a><br>
<a href="./call.go#L6"><code>(1 2;0)[0][1] /2</code></a><br>
<a href="./call.go#L11"><code>3+ /3+</code>(projection)</a><br>
<a href="./call.go#L12"><code>+[3;] /3+</code></a><br>
<a href="./call.go#L13"><code>+[;3] /+[;3]</code></a><br>
<a href="./call.go#L14"><code>3 imag /3 imag</code></a><br>
<a href="./call.go#L15"><code>+ /+</code></a><br>
<a href="./call.go#L16"><code>-[1] /-1</code></a><br>
<a href="./call.go#L17"><code>+/[1 2 3] /6</code></a><br>
<a href="./call.go#L18"><code>-[5;3] /2</code></a><br>
<a href="./call.go#L67"><code>{x+y}[1;2] /3</code></a><br>
<a href="./call.go#L68"><code>{x+y}.1 2 /3</code></a><br>
<a href="./index.go#L63"><code>(1 2 3).(1) /2</code></a><br>
<a href="./index.go#L64"><code>(1 2 3;4 5 6).(1 2) /6</code></a><br>
<a href="./index.go#L65"><code>(1 2 3;4 5 6)[1;2] /6</code></a><br>
<a href="./index.go#L66"><code>(1 2 3;4 5 6)[1 0;1] /5 2</code></a><br>
<a href="./index.go#L67"><code>(1 2 3;4 5 6)[0 1;1 0] /(2 1;5 4)</code></a><br>
<a href="./verbs.go#L282"><code>.&#34;1+2&#34; /3</code></a><br>
</details>
<details><summary><code>abs</code>abs </summary>
<a href="./numeric.go#L365"><code>abs -1 2 /1 2</code></a><br>
<a href="./numeric.go#L366"><code>abs 2a30 /2.</code></a><br>
</details>
<details><summary><code>angle</code>angle rotate</summary>
<a href="./numeric.go#L404"><code>angle 1a20 2a45 /20 45.</code></a><br>
<a href="./numeric.go#L405"><code>angle 1.2 /0.</code></a><br>
<a href="./numeric.go#L433"><code>1a20 angle 25 /1a45</code></a><br>
</details>
<details><summary><code>real</code>zreal </summary>
<a href="./numeric.go#L466"><code>real 1a300 /0.5</code></a><br>
</details>
<details><summary><code>imag</code>zimag complx</summary>
<a href="./numeric.go#L482"><code>imag 1a60 /0.8660254037844386</code></a><br>
<a href="./numeric.go#L498"><code>1 imag 1 /1.4142135623730951a45</code></a><br>
</details>
<details><summary><code>conj</code>conj </summary>
<a href="./numeric.go#L501"><code>conj 1a60 /1a300</code></a><br>
</details>
<details><summary><code>read</code>read readdir</summary>
<a href="./read.go#L10"><code>@read&#34;readme.md&#34; /`C</code>(read file)</a><br>
<a href="./read.go#L11"><code>#read(/&#34;\\.md$&#34;) /1</code>(filter cwd)</a><br>
<a href="./read.go#L34"><code>@(/&#34;\\.md$&#34;)read` /`L</code>(filter dir)</a><br>
</details>
<details><summary><code>csv</code>csvwrite csvread</summary>
<a href="./csv.go#L13"><code>&#34;&#34;csv&#34;ab|cd|ef\ngh|ij|kl\n&#34; /((&#34;ab&#34;;&#34;gh&#34;);(&#34;cd&#34;;&#34;ij&#34;);(&#34;ef&#34;;&#34;kl&#34;))</code>(auto-detect)</a><br>
<a href="./csv.go#L14"><code>&#34;,if&#34;csv&#34;1,2\n3,4\n5,6\n&#34; /(1 3 5;2 4 6.)</code></a><br>
<a href="./csv.go#L15"><code>&#34;,2hiffs&#34;csv&#34;x\n\n1,2,0,abc\n2,3,90,gh&#34; /(1 2;2 3.;0 90.;`abc`gh)</code></a><br>
<a href="./csv.go#L16"><code>&#34;;izs&#34;csv&#34;1;2;0;abc\n2;3;90;gh&#34; /(1 2;2a0 3a90;`abc`gh)</code></a><br>
</details>
<details><summary><code>solve</code>qr solve</summary>
</details>
<details><summary><code>&#39;</code>each each2</summary>
<a href="./adverbs.go#L31"><code>-&#39;1 2 3 /-1 -2 -3</code></a><br>
<a href="./adverbs.go#L49"><code>(1;2 3)+&#39;(2;4 5) /(3;6 8)</code></a><br>
<a href="./adverbs.go#L50"><code>(1;2 3)+&#39;5 /(6;7 8)</code></a><br>
</details>
<details><summary><code>/</code>over over2</summary>
<a href="./adverbs.go#L79"><code>+/1 2 3 /6</code></a><br>
<a href="./adverbs.go#L80"><code>(+)/1 2 3 /6</code></a><br>
<a href="./find.go#L403"><code>&#34;n&#34;/(&#34;ab&#34;;&#34;cde&#34;) /&#34;abncde&#34;</code></a><br>
<a href="./find.go#L404"><code>1 2/(!3;!2) /0 1 2 1 2 0 1</code></a><br>
<a href="./find.go#L405"><code>(1+2)/(1 2;7 8) /1 2 3 7 8</code></a><br>
<a href="./find.go#L469"><code>&#34;pq&#34;(&#34;ab&#34;)/&#34;alaba&#34; /&#34;alpqa&#34;</code>(replace)</a><br>
<a href="./find.go#L470"><code>&#34;bb&#34;(&#34;a&#34;)/&#34;alpha&#34; /&#34;bblphbb&#34;</code></a><br>
<a href="./find.go#L471"><code>&#34;$2$1&#34;(/&#34;.([0-9])..([0-9]).&#34;)/&#34;12345678&#34; /&#34;5278&#34;</code>(regex-replace)</a><br>
</details>
<details><summary><code>\</code>scan scan2</summary>
<a href="./adverbs.go#L132"><code>-\1 2 3 /1 -1 -4</code></a><br>
<a href="./find.go#L359"><code>&#34;x&#34;\&#34;abxdexxg&#34; /(&#34;ab&#34;;&#34;de&#34;;&#34;&#34;;,&#34;g&#34;)</code>(split)</a><br>
<a href="./find.go#L360"><code>&#34;&#34;\&#34;ab  de f &#34; /(&#34;ab&#34;;&#34;de&#34;;,&#34;f&#34;)</code>(fields)</a><br>
<a href="./find.go#L361"><code>(/&#34;[0-9]&#34;)\&#34;ab3cd2cv4&#34; /(&#34;ab&#34;;&#34;cd&#34;;&#34;cv&#34;;&#34;&#34;)</code>(regexp-split)</a><br>
</details>
<details><summary><code>&#39;:</code>pairs pairs2</summary>
<a href="./adverbs.go#L216"><code>&lt;&#39;:4 2 1 0 2 /01110b</code></a><br>
</details>
<details><summary><code>/:</code>fix eachright</summary>
</details>
<details><summary><code>\:</code>scanfix eachleft</summary>
</details>
<details><summary><code>:</code> assign</summary>
<a href="./assign.go#L27"><code>x:3 /3</code></a><br>
<a href="./assign.go#L28"><code>(x;y):1 /1</code></a><br>
<a href="./assign.go#L29"><code>(x;y):1 2 /1 2</code>(destructing assign)</a><br>
<a href="./assign.go#L30"><code>(x;y):1 2 3 /1 2 3</code></a><br>
<a href="./assign.go#L31"><code>z:(x;y):1 2 3 /1 2 3</code></a><br>
<a href="./assign.go#L32"><code>a+:a:3 /6</code>(modified assign)</a><br>
<a href="./assign.go#L33"><code>a:!5;a[!2]:0 /0 0 2 3 4</code>(indexed assign)</a><br>
<a href="./assign.go#L34"><code>x:!5;x[2 3]+:1 /0 1 3 4 4</code>(indexed modified)</a><br>
<a href="./assign.go#L35"><code>x:2^!6;x[1;2]:9 /(0 1 2;3 4 9)</code>(at-depth)</a><br>
<a href="./assign.go#L36"><code>x:2^!6;x[!2;!2]:9 /(9 9 2;9 9 5)</code>(matrix-assign)</a><br>
<a href="./assign.go#L37"><code>x:2^!6;x[;1]*:10 /(0 10 2;3 40 5)</code>(column)</a><br>
</details>
<details><summary><code>λ</code> func</summary>
<a href="./lambda.go#L3"><code>{1+2} /{1+2}</code></a><br>
<a href="./lambda.go#L4"><code>{1+2}[] /3</code></a><br>
<a href="./lambda.go#L5"><code>{x+2}[2] /4</code></a><br>
<a href="./lambda.go#L6"><code>{x+y}[2;3] /5</code></a><br>
<a href="./lambda.go#L7"><code>{[a;b]3*a+b}[3;4] /21</code></a><br>
<a href="./lambda.go#L8"><code>{[]3}[] /3</code></a><br>
<a href="./lambda.go#L9"><code>{(a;y*a:x)}[2;3] /2 6</code></a><br>
<a href="./lambda.go#L10"><code>{a+y*a:x}[2;3] /8</code></a><br>
<a href="./lambda.go#L11"><code>a+{2*a:x}[2]+a:1 /6</code>(local assign)</a><br>
<a href="./lambda.go#L12"><code>a+{2*a::x}[2]+a:1 /7</code>(global assign)</a><br>
</details>
<details><summary><code>;</code> statement</summary>
<a href="./exec.go#L163"><code>1;2 /2</code></a><br>
</details>
