(() => {
// this is a modified version of tinygo/targets/wasm_exe.js 
// it writes error messages (panic output) to the kons textarea instead of console.log
	if (typeof window !== "undefined") { window.global = window; } else if (typeof self !== "undefined") { self.global = self; } else { throw new Error("cannot export Go (neither window nor self is defined)"); }
	global.pnk = "";
	const decoder = new TextDecoder("utf-8");
	var logLine = [];
	global.Go = class {
		constructor() {
			const mem = () => { return new DataView(this._inst.exports.memory.buffer); }
			this.importObject = {
				env: {
					io_get_stdout: function() { return 1; },
					resource_write: function(fd, ptr, len) {
						if (fd == 1) {
							for (let i=0; i<len; i++) {
								let c = mem().getUint8(ptr+i);
								if (c == 13) { // CR // ignore } else if (c == 10) { // LF
									// write line
									let line = decoder.decode(new Uint8Array(logLine));
									logLine = [];
									// console.log(line);
									pnk = line // read by try/catch on panics
								} else { logLine.push(c); }
							}
						} else { console.error('invalid file descriptor:', fd); }
					},
				}
			};
		}
		async run(instance) {
			this._inst = instance; this._values = [ NaN, 0, null, true, false, global, this._inst.exports.memory, this ]; this._refs = new Map(); this._callbackShutdown = false; this.exited = false;
			const mem = new DataView(this._inst.exports.memory.buffer);
			while (true) {
				const callbackPromise = new Promise((resolve) => { this._resolveCallbackPromise = () => { if (this.exited) { throw new Error("bad callback: Go program has already exited"); } setTimeout(resolve, 0); }; });
				this._inst.exports.cwa_main(); if (this.exited) { break; }
				await callbackPromise;
			}
		}
		_resume() {
			if (this.exited) { throw new Error("Go program has already exited"); }
			this._inst.exports.resume();
			if (this.exited) { this._resolveExitPromise(); }
		}
		_makeFuncWrapper(id) { const go = this; return function () { const event = { id: id, this: this, args: arguments }; go._pendingEvent = event; go._resume(); return event.result; }; }
	}
})();
