(window.webpackJsonp = window.webpackJsonp || []).push([["app"], {
    0: function (t, n, e) {
        t.exports = e("56d7")
    }, "034f": function (t, n, e) {
        "use strict";
        e("85ec")
    }, 4362: function (t, n, e) {
        var r, i;
        n.nextTick = function (t) {
            var n = Array.prototype.slice.call(arguments);
            n.shift(), setTimeout((function () {
                t.apply(null, n)
            }), 0)
        }, n.platform = n.arch = n.execPath = n.title = "browser", n.pid = 1, n.browser = !0, n.env = {}, n.argv = [], n.binding = function (t) {
            throw new Error("No such module. (Possibly not yet loaded)")
        }, i = "/", n.cwd = function () {
            return i
        }, n.chdir = function (t) {
            r || (r = e("df7c")), i = r.resolve(t, i)
        }, n.exit = n.kill = n.umask = n.dlopen = n.uptime = n.memoryUsage = n.uvCounters = function () {
        }, n.features = {}
    }, "56d7": function (t, n, e) {
        "use strict";
        e.r(n);
        e("e260"), e("e6cf"), e("cca6"), e("a79d");
        var r = e("2b0e"), i = (e("b0c0"), e("bc3a")), o = e.n(i), a = {
            name: "App", created: function () {
                this.flush()
            }, data: function () {
                return {
                    url: "/api",
                    indexCnt: 0,
                    totalCnt: 0,
                    curImage: null,
                    matchImageList: null,
                    onLoad: !1,
                    info: {tp5_ac: null, tp5_rc: null, tp10_ac: null, tp10_rc: null, p50_ac: null, p50_rc: null}
                }
            }, components: {}, methods: {
                flush: function () {
                    var t = this;
                    o.a.get(this.url + "/num").then((function (n) {
                        return t.indexCnt = n.data.data
                    })).catch((function (t) {
                        return console.log(t)
                    })), o.a.get(this.url + "/tot_num").then((function (n) {
                        return t.totalCnt = n.data.data
                    })).catch((function (t) {
                        return console.log(t)
                    }))
                }, upload: function () {
                    for (var t = this, n = this.$refs.up.files, e = {headers: {"Content-Type": "multipart/form-data"}}, r = 0; r < n.length; ++r) {
                        var i = new FormData;
                        i.append("file", n[r]), o.a.post(this.url + "/upload", i, e).then((function () {
                            return t.flush()
                        })).catch((function (t) {
                            return console.log(t)
                        }))
                    }
                }, find: function () {
                    var t = this, n = this.$refs.fi.files[0], e = new FormData;
                    e.append("file", n), e.append("file", n), o.a.post(this.url + "/find", e, {headers: {"Content-Type": "multipart/form-data"}}).then((function (e) {
                        t.curImage = n.name, t.info = e.data.data.info, t.matchImageList = e.data.data.data, console.log(t.matchImageList)
                    })).catch((function (t) {
                        return console.log(t)
                    }))
                }, reduce: function () {
                    var t = this;
                    this.onLoad = !0, o.a.post(this.url + "/reduce", {word: 64}).then((function () {
                        t.flush(), t.onLoad = !1
                    }))
                }
            }
        };
        e("034f");
        var s = function (t, n, e, r, i, o, a, s) {
            var c, l = "function" == typeof t ? t.options : t;
            if (n && (l.render = n, l.staticRenderFns = e, l._compiled = !0), r && (l.functional = !0), o && (l._scopeId = "data-v-" + o), a ? (c = function (t) {
                (t = t || this.$vnode && this.$vnode.ssrContext || this.parent && this.parent.$vnode && this.parent.$vnode.ssrContext) || "undefined" == typeof __VUE_SSR_CONTEXT__ || (t = __VUE_SSR_CONTEXT__), i && i.call(this, t), t && t._registeredComponents && t._registeredComponents.add(a)
            }, l._ssrRegister = c) : i && (c = s ? function () {
                i.call(this, (l.functional ? this.parent : this).$root.$options.shadowRoot)
            } : i), c) if (l.functional) {
                l._injectStyles = c;
                var h = l.render;
                l.render = function (t, n) {
                    return c.call(n), h(t, n)
                }
            } else {
                var u = l.beforeCreate;
                l.beforeCreate = u ? [].concat(u, c) : [c]
            }
            return {exports: t, options: l}
        }(a, (function () {
            var t = this, n = t.$createElement, e = t._self._c || n;
            return e("div", {attrs: {id: "app"}}, [e("div", [e("h1", [t._v("上传图片")]), e("div", {staticClass: "upload"}, [e("svg", {
                staticClass: "icon",
                attrs: {
                    t: "1620918057417",
                    viewBox: "0 0 1024 1024",
                    version: "1.1",
                    xmlns: "http://www.w3.org/2000/svg",
                    width: "100",
                    height: "100"
                }
            }, [e("path", {
                attrs: {
                    d: "M1035.3664 1035.82151111h-86.47111111v-27.30666666h59.16444444v-59.16444445h27.30666667zM881.69813333 1035.82151111H747.30382222v-27.30666666h134.39431111v27.30666666z m-201.59146666 0H545.71235555v-27.30666666h134.39431112v27.30666666z m-201.59146667 0H344.12088889v-27.30666666h134.39431111v27.30666666z m-201.59146667 0H142.52942222v-27.30666666h134.39431111v27.30666666zM75.33226667 1035.82151111h-86.47111112v-86.47111111h27.30666667v59.16444445h59.16444445zM16.16782222 882.15324445h-27.30666667V747.75893333h27.30666667v134.39431112z m0-201.59146667h-27.30666667V546.16746667h27.30666667v134.39431111z m0-201.59146667h-27.30666667V344.576h27.30666667v134.39431111z m0-201.59146666h-27.30666667V142.98453333h27.30666667v134.39431112zM16.16782222 75.78737778h-27.30666667v-86.47111111h86.47111112v27.30666666h-59.16444445zM881.69813333 16.62293333H747.30382222v-27.30666666h134.39431111v27.30666666z m-201.59146666 0H545.71235555v-27.30666666h134.39431112v27.30666666z m-201.59146667 0H344.12088889v-27.30666666h134.39431111v27.30666666z m-201.59146667 0H142.52942222v-27.30666666h134.39431111v27.30666666zM1035.3664 75.78737778h-27.30666667v-59.16444445h-59.16444444v-27.30666666h86.47111111zM1035.3664 882.15324445h-27.30666667V747.75893333h27.30666667v134.39431112z m0-201.59146667h-27.30666667V546.16746667h27.30666667v134.39431111z m0-201.59146667h-27.30666667V344.576h27.30666667v134.39431111z m0-201.59146666h-27.30666667V142.98453333h27.30666667v134.39431112z",
                    fill: "#bfbfbf"
                }
            }), e("path", {
                attrs: {
                    d: "M599.74674456 523.1642475H424.34544577c-6.10760857 0-11.06283816-4.95062007-11.06283815-11.06283815 0-6.10760857 4.95522958-11.06283816 11.06283814-11.06283815h175.4012988c6.10760857 0 11.06283816 4.95522958 11.06283816 11.06283815 0 6.11221808-4.95062007 11.06283816-11.06283816 11.06283815z",
                    fill: "#bfbfbf"
                }
            }), e("path", {
                attrs: {
                    d: "M512.04609516 610.86489689c-6.10760857 0-11.06283816-4.95522958-11.06283816-11.06283815V424.40075995c0-6.10760857 4.95522958-11.06283816 11.06283816-11.06283814s11.06283816 4.95522958 11.06283814 11.06283814v175.40129879c0 6.11221808-4.95522958 11.06283816-11.06283814 11.06283815z",
                    fill: "#bfbfbf"
                }
            })]), e("input", {
                ref: "up",
                staticClass: "fileInput",
                attrs: {type: "file", multiple: ""},
                on: {change: t.upload}
            })]), e("p", [t._v(" 当前索引中一共有 " + t._s(t.indexCnt) + " 张图片 "), e("span", [t._v(" , ")]), t._v(" 当前图库中一共有 " + t._s(t.totalCnt) + " 张图片 ")])]), e("button", {
                attrs: {disabled: t.onLoad},
                on: {click: t.reduce}
            }, [t._v("处理所有上传的图片")]), e("hr"), e("div", [e("h1", [t._v("检索图片")]), e("div", {staticClass: "upload"}, [e("svg", {
                staticClass: "icon",
                attrs: {
                    t: "1620918057417",
                    viewBox: "0 0 1024 1024",
                    version: "1.1",
                    xmlns: "http://www.w3.org/2000/svg",
                    width: "100",
                    height: "100"
                }
            }, [e("path", {
                attrs: {
                    d: "M1035.3664 1035.82151111h-86.47111111v-27.30666666h59.16444444v-59.16444445h27.30666667zM881.69813333 1035.82151111H747.30382222v-27.30666666h134.39431111v27.30666666z m-201.59146666 0H545.71235555v-27.30666666h134.39431112v27.30666666z m-201.59146667 0H344.12088889v-27.30666666h134.39431111v27.30666666z m-201.59146667 0H142.52942222v-27.30666666h134.39431111v27.30666666zM75.33226667 1035.82151111h-86.47111112v-86.47111111h27.30666667v59.16444445h59.16444445zM16.16782222 882.15324445h-27.30666667V747.75893333h27.30666667v134.39431112z m0-201.59146667h-27.30666667V546.16746667h27.30666667v134.39431111z m0-201.59146667h-27.30666667V344.576h27.30666667v134.39431111z m0-201.59146666h-27.30666667V142.98453333h27.30666667v134.39431112zM16.16782222 75.78737778h-27.30666667v-86.47111111h86.47111112v27.30666666h-59.16444445zM881.69813333 16.62293333H747.30382222v-27.30666666h134.39431111v27.30666666z m-201.59146666 0H545.71235555v-27.30666666h134.39431112v27.30666666z m-201.59146667 0H344.12088889v-27.30666666h134.39431111v27.30666666z m-201.59146667 0H142.52942222v-27.30666666h134.39431111v27.30666666zM1035.3664 75.78737778h-27.30666667v-59.16444445h-59.16444444v-27.30666666h86.47111111zM1035.3664 882.15324445h-27.30666667V747.75893333h27.30666667v134.39431112z m0-201.59146667h-27.30666667V546.16746667h27.30666667v134.39431111z m0-201.59146667h-27.30666667V344.576h27.30666667v134.39431111z m0-201.59146666h-27.30666667V142.98453333h27.30666667v134.39431112z",
                    fill: "#bfbfbf"
                }
            }), e("path", {
                attrs: {
                    d: "M599.74674456 523.1642475H424.34544577c-6.10760857 0-11.06283816-4.95062007-11.06283815-11.06283815 0-6.10760857 4.95522958-11.06283816 11.06283814-11.06283815h175.4012988c6.10760857 0 11.06283816 4.95522958 11.06283816 11.06283815 0 6.11221808-4.95062007 11.06283816-11.06283816 11.06283815z",
                    fill: "#bfbfbf"
                }
            }), e("path", {
                attrs: {
                    d: "M512.04609516 610.86489689c-6.10760857 0-11.06283816-4.95522958-11.06283816-11.06283815V424.40075995c0-6.10760857 4.95522958-11.06283816 11.06283816-11.06283814s11.06283816 4.95522958 11.06283814 11.06283814v175.40129879c0 6.11221808-4.95522958 11.06283816-11.06283814 11.06283815z",
                    fill: "#bfbfbf"
                }
            })]), e("input", {
                ref: "fi",
                staticClass: "fileInput",
                attrs: {type: "file"},
                on: {change: t.find}
            })]), e("hr"), null !== t.curImage ? e("div", {
                staticStyle: {
                    display: "inline-block",
                    "border-style": "dotted",
                    "border-width": "2px",
                    "border-color": "yellowgreen",
                    padding: "2px 2px"
                }
            }, [t._v(" 当前图片 "), e("br"), e("img", {
                staticStyle: {"max-width": "300px", "max-height": "300px"},
                attrs: {src: "http://localhost:5000/api/img/" + t.curImage, alt: "cur"}
            })]) : t._e(), e("br"), null !== t.matchImageList ? e("div", [e("p", [t._v(" tp5 准确度: " + t._s(t.info.tp5_ac) + " "), e("span", [t._v(",")]), t._v(" 召回率: " + t._s(t.info.tp5_rc) + " "), e("br"), t._v(" tp10 准确度: " + t._s(t.info.tp10_ac) + " "), e("span", [t._v(",")]), t._v(" 召回率: " + t._s(t.info.tp10_rc) + " "), e("br"), t._v(" 匹配度 > 0.5 的准确度: " + t._s(t.info.p50_ac) + " "), e("span", [t._v(",")]), t._v(" 召回率: " + t._s(t.info.p50_rc) + " "), e("br")]), t._l(t.matchImageList, (function (n) {
                return e("div", {
                    key: n,
                    staticStyle: {
                        display: "inline-block",
                        "border-style": "dotted",
                        "border-width": "2px",
                        "border-color": "burlywood",
                        padding: "2px 2px"
                    }
                }, [t._v(" 相似度：" + t._s(n.value) + " "), e("br"), null !== t.curImage ? e("img", {
                    staticStyle: {
                        "max-width": "300px",
                        "max-height": "300px"
                    }, attrs: {src: "http://localhost:5000/api/img/" + n.name, alt: "cur"}
                }) : t._e()])
            }))], 2) : t._e()])])
        }), [], !1, null, null, null).exports;
        r.a.config.productionTip = !1, new r.a({
            render: function (t) {
                return t(s)
            }
        }).$mount("#app")
    }, "85ec": function (t, n, e) {
    }, c8ba: function (t, n) {
        var e;
        e = function () {
            return this
        }();
        try {
            e = e || new Function("return this")()
        } catch (t) {
            "object" == typeof window && (e = window)
        }
        t.exports = e
    }, df7c: function (t, n, e) {
        (function (t) {
            function e(t, n) {
                for (var e = 0, r = t.length - 1; r >= 0; r--) {
                    var i = t[r];
                    "." === i ? t.splice(r, 1) : ".." === i ? (t.splice(r, 1), e++) : e && (t.splice(r, 1), e--)
                }
                if (n) for (; e--; e) t.unshift("..");
                return t
            }

            function r(t, n) {
                if (t.filter) return t.filter(n);
                for (var e = [], r = 0; r < t.length; r++) n(t[r], r, t) && e.push(t[r]);
                return e
            }

            n.resolve = function () {
                for (var n = "", i = !1, o = arguments.length - 1; o >= -1 && !i; o--) {
                    var a = o >= 0 ? arguments[o] : t.cwd();
                    if ("string" != typeof a) throw new TypeError("Arguments to path.resolve must be strings");
                    a && (n = a + "/" + n, i = "/" === a.charAt(0))
                }
                return (i ? "/" : "") + (n = e(r(n.split("/"), (function (t) {
                    return !!t
                })), !i).join("/")) || "."
            }, n.normalize = function (t) {
                var o = n.isAbsolute(t), a = "/" === i(t, -1);
                return (t = e(r(t.split("/"), (function (t) {
                    return !!t
                })), !o).join("/")) || o || (t = "."), t && a && (t += "/"), (o ? "/" : "") + t
            }, n.isAbsolute = function (t) {
                return "/" === t.charAt(0)
            }, n.join = function () {
                var t = Array.prototype.slice.call(arguments, 0);
                return n.normalize(r(t, (function (t, n) {
                    if ("string" != typeof t) throw new TypeError("Arguments to path.join must be strings");
                    return t
                })).join("/"))
            }, n.relative = function (t, e) {
                function r(t) {
                    for (var n = 0; n < t.length && "" === t[n]; n++) ;
                    for (var e = t.length - 1; e >= 0 && "" === t[e]; e--) ;
                    return n > e ? [] : t.slice(n, e - n + 1)
                }

                t = n.resolve(t).substr(1), e = n.resolve(e).substr(1);
                for (var i = r(t.split("/")), o = r(e.split("/")), a = Math.min(i.length, o.length), s = a, c = 0; c < a; c++) if (i[c] !== o[c]) {
                    s = c;
                    break
                }
                var l = [];
                for (c = s; c < i.length; c++) l.push("..");
                return (l = l.concat(o.slice(s))).join("/")
            }, n.sep = "/", n.delimiter = ":", n.dirname = function (t) {
                if ("string" != typeof t && (t += ""), 0 === t.length) return ".";
                for (var n = t.charCodeAt(0), e = 47 === n, r = -1, i = !0, o = t.length - 1; o >= 1; --o) if (47 === (n = t.charCodeAt(o))) {
                    if (!i) {
                        r = o;
                        break
                    }
                } else i = !1;
                return -1 === r ? e ? "/" : "." : e && 1 === r ? "/" : t.slice(0, r)
            }, n.basename = function (t, n) {
                var e = function (t) {
                    "string" != typeof t && (t += "");
                    var n, e = 0, r = -1, i = !0;
                    for (n = t.length - 1; n >= 0; --n) if (47 === t.charCodeAt(n)) {
                        if (!i) {
                            e = n + 1;
                            break
                        }
                    } else -1 === r && (i = !1, r = n + 1);
                    return -1 === r ? "" : t.slice(e, r)
                }(t);
                return n && e.substr(-1 * n.length) === n && (e = e.substr(0, e.length - n.length)), e
            }, n.extname = function (t) {
                "string" != typeof t && (t += "");
                for (var n = -1, e = 0, r = -1, i = !0, o = 0, a = t.length - 1; a >= 0; --a) {
                    var s = t.charCodeAt(a);
                    if (47 !== s) -1 === r && (i = !1, r = a + 1), 46 === s ? -1 === n ? n = a : 1 !== o && (o = 1) : -1 !== n && (o = -1); else if (!i) {
                        e = a + 1;
                        break
                    }
                }
                return -1 === n || -1 === r || 0 === o || 1 === o && n === r - 1 && n === e + 1 ? "" : t.slice(n, r)
            };
            var i = "b" === "ab".substr(-1) ? function (t, n, e) {
                return t.substr(n, e)
            } : function (t, n, e) {
                return n < 0 && (n = t.length + n), t.substr(n, e)
            }
        }).call(this, e("4362"))
    }
}, [[0, "runtime", "npm.core-js", "npm.axios", "npm.vue"]]]);
//# sourceMappingURL=app.3b456eb4.js.map