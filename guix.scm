(use-modules (guix)
             (guix build-system copy)
             ((guix licenses) #:prefix license:)
            (nonguix licenses)
             (gnu packages golang)
             (gnu packages debug)
             (gnu packages embedded)
             (gnu packages bash)
             (gnu packages shells)
             (gnu packages linux)
             (th packages helix-editor)
             (th packages starship)(gnu packages python)(gnu packages perl)
            (nonguix build-system binary)
             )


; (define fish-container
;   (package
;     (inherit fish)
;    (arguments
;      '(#:tests? #f
;        #:phases (modify-phases %standard-phases
;           (delete 'patch-fish-config))))                          ;funky version number
; ))

; (define tinygo-bin
;   (package
;     (name "tinygo-bin")
;     (version "0.31.2")
;     (source
;      (origin
;        (method url-fetch)
;        (uri (string-append
;              "https://github.com/tinygo-org/tinygo/releases/download/v"
;              version
;              "/tinygo"
; 	     version
; 	     ".linux-amd64.tar.gz"))
;        (sha256
;         (base32
;          "1f7sv37n6av94vgv2rbpy2mqvsirn2hapj6a5r6i46g6xzlkma28"))))
;     (build-system binary-build-system)
;     ; (inputs (list python perl))
;     (arguments
;      `(#:install-plan
;       `(("./" "/"))
;       ; ("src/" "/src")
;       ; ("lib/" "/lib")
;       ; ("pkg/" "/pkg"))
; 		 ))
;     (synopsis "TinyGo - Go compiler for small places")
;     (description "Go compiler for small places. Microcontrollers, WebAssembly (WASM/WASI), and command-line tools. Based on LLVM.")
;     (home-page "https://tinygo.org/")
;     (license (nonfree "https://tinygo.org"))
; ))

(define dev-env
  (package
    (name "dev-env")
    (version "0.1")
    (license #f)
    (source #f)
    (description "")
    (home-page "")
    (synopsis "")
    (build-system copy-build-system)
   (arguments
     '(#:tests? #f
       #:phases (modify-phases %standard-phases
          (delete 'unpack))))                          ;funky version number
    (propagated-inputs
      (list go-1.21 delve gopls go-github-com-prometheus-client-golang))
        ; ("openocd" ,openocd)
        ; ("go" ,go)
        ; ("delve" ,delve)
        ; ("gopls" ,gopls)
        ; ("starship-bin" ,starship-bin)
        ; ("helix-editor-bin" ,helix-editor-bin)
        ; ("tinygo-bin" ,tinygo-bin)))
))
dev-env
