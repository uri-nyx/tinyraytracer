#lang racket
(require math/array)
(require racket/vector)

; step 1, create an rgb image
; The format is PPM. We need a header and a body of pixels
; The header is comprised of a format identifier ("P6")
; a width, height, and te max value for each channel
(define (ppm-header width height) (string->bytes/latin-1 (string-join (list "P6" (number->string width) (number->string height) "255" "\n") " ")))

; The pixel buffer is an array of pixels of length w*h
; a pixel is an array of three bytes, one for each channel
(define (make-pixels width height) (array->mutable-array (make-array (vector (* width height)) #(0 0 0))))
(define (array->bytes pixels) (list->bytes (flatten (map vector->list (array->list pixels)))))

(define (save-ppm filename header pixels)
  (with-output-to-file filename
    (lambda () (write-bytes (bytes-join (list header pixels) #"")))
    #:mode 'binary #:exists 'replace))

(define (test-pixels w h)(
                          (letrec ([header (ppm-header w h)]
                                   [pixels (make-pixels w h)])
                           (set-normalized-color! pixels #(0) 1 0 0)
                           (set-normalized-color! pixels #(1) 0 1 0)
                           (set-normalized-color! pixels #(2) 0 0 1)
                           (set-normalized-color! pixels (vector (sub1 (* w h))) 1 1 1)

                           (lambda () (save-ppm "test.ppm" header (array->bytes pixels))))))

; Very imperative, don't judge me please. I don't know why it is reversed and mirrored
(define (test-ppm w h)(
                       (letrec ([header (ppm-header w h)]
                                [pixels (make-pixels w h)])

                         (for ([j (in-list (build-list h values))])
                            (printf "scanlines remaining: ~a\n" j)
                           (for ([i (in-list (build-list w values))])
                             (let ([r (/ i (- w 1))] [g (/ j (- h 1))] [b 0.25])
                              (set-normalized-color! pixels (vector (+ (* (- j (sub1 h)) (- w)) i)) r g b))))
                        (display "Done, now saving to file...")
                        (lambda () (save-ppm "test.ppm" header (array->bytes pixels))))))

; Vector (pixel) operations
(define (v-x vec) (vector-ref vec 0))
(define (v-y vec) (vector-ref vec 1))
(define (v-z vec) (vector-ref vec 2))
(define (v-set-x! vec n) (vector-set! vec 0 n))
(define (v-set-y! vec n) (vector-set! vec 1 n))
(define (v-set-z! vec n) (vector-set! vec 2 n))

; Per component operations
(define (v-aritmetic proc) (curry vector-map proc))
(define (v-mul-scalar vec scalar) (vector-map (curry * scalar) vec))
(define (v-div-scalar vec scalar) (vector-map (curry * (/ 1 scalar)) vec))
(define (v-expt power) (curry vector-map (curryr expt power)))
(define v-add (v-aritmetic +))
(define v-sub (v-aritmetic -))


(define (v-length-squared vec) (apply + (vector->list ((v-expt 2) vec))))
(define (v-length vec) (sqrt (v-length-squared vec))) ; sqrt(v1^2 + v2^2 + v3^2)
(define (v-unit vec) (v-div-scalar vec (v-length vec))) ; v / sqrt(v1^2 + v2^2 + v3^2)

; Products
(define v-dot-prod (compose (curry apply +) vector->list (curry vector-map *)))
(define (v3-cross-prod vec1 vec2)
  (match-let ([(vector v1x v1y v1z) vec1] [(vector v2x v2y v2z) vec2])
    (vector (- (* v1y v2z) (* v1z v2y))
            (- (* v1z v2x) (* v1x v2z))
            (- (* v1x v2y) (* v1y v2x)))))

; Assignment of values to pixels
(define (get-index w col row) (vector (+ (* w row) col)))
(define (set-color! buffer index r g b) (array-set! buffer index (vector r g b)))

(define (set-normalized-color! buffer index r g b)
  (let ([norm->byte (lambda (x) (exact-floor (* x 255.999)))])
    (array-set! buffer index (vector (norm->byte r) (norm->byte g) (norm->byte b)))))

(define (set-normalized-color-v! buffer index color)
  (let ([v-norm->byte (lambda (v) (vector-map exact-floor (v-mul-scalar v 255.999)))])
    (array-set! buffer index (v-norm->byte color))))

; sytactic sugar
(define color vector)
(define point vector)

; Rays have direction (vector) and an origin (point)
(struct ray (origin direction))
(define (ray-at ray scalar) (v-add (ray-origin ray) (v-mul-scalar (ray-direction ray) scalar)))

(define (ray-color ray) (letrec ([unit-direction (v-unit (ray-direction ray))]
                                 [s (sphere (point 0 0 -1) 0.5)]
                                 [t (hit-sphere? s ray)])
                                
                                (if (> t 0)
                                  (let ([normal (v-unit (v-sub (ray-at ray t) (vector 0 0 -1)))])
                                   (v-mul-scalar (vector-map add1 normal) 0.5))
                                ; else
                                  (let ([g (* 0.5 (+ 1 (v-y unit-direction)))])
                                   (v-add
                                    (v-mul-scalar (color 1 1 1) (- 1 g))
                                    (v-mul-scalar (color 0.5 0.7 1) g))))))

(struct image (aspect-ratio width height))
(struct camera (view-height view-width focal origin horizontal vertical lower-left-corner))

(define (make-image aspect-ratio width) (image aspect-ratio width (exact-floor (/ width aspect-ratio))))
(define (make-camera image view-height focal origin) (letrec ([view-width (* view-height (image-aspect-ratio image))]
                                                              [horizontal (vector view-width 0 0)]
                                                              [vertical (vector 0 view-height 0)]
                                                              [lo-left-corner (v-sub origin (v-div-scalar horizontal 2)
                                                                               (v-div-scalar vertical 2)
                                                                               (vector 0 0 focal))])
                                                          (camera view-height view-width focal origin horizontal vertical lo-left-corner)))

(define (render image camera) (letrec ([image-w (image-width image)]
                                       [image-h (image-height image)]
                                       [header (ppm-header image-w image-h)]
                                       [pixels (make-pixels image-w image-h)])
                                
                                (displayln (list "Rendering image..." image-w image-h))
                                       
                                (for ([j (in-list (build-list image-h values))])
                                  (printf "scanlines remaining: ~a\n" (- image-h j))

                                  (for ([i (in-list (build-list image-w values))])

                                    (letrec ([u (/ i (sub1 image-w))]
                                             [v (/ j (sub1 image-h))]
                                             [r (ray (camera-origin camera)
                                                     (v-sub (v-add (camera-lower-left-corner camera) 
                                                               (v-mul-scalar (camera-horizontal camera) u)
                                                               (v-mul-scalar (camera-vertical camera) v))
                                                         (camera-origin camera)))]
                                              ;lower_left_corner + u*horizontal + v*vertical - origin)
                                             [idx (get-index (- image-w) i (- j (sub1 image-h)))] ;(vector (+ (* (- j (sub1 h)) (- w)) i)
                                             [color (ray-color r)])
                                          
                                      (set-normalized-color! pixels idx (v-x color) (v-y color) (v-z color)))))

                                (display "Done, now saving to file...")
                                (save-ppm "test.ppm" header (array->bytes pixels))))

(define (test-render) (letrec ([image (make-image 16/9 400)] [camera (make-camera image 2 1 (point 0 0 0))])
                        (render image camera)))

(struct sphere (center radius))

(define (hit-sphere? sphere ray)
  ;quaderatic equation
 (letrec ([oc (v-sub (ray-origin ray) (sphere-center sphere))]
          [a (v-length-squared (ray-direction ray))]
          [b/2 (v-dot-prod oc (ray-direction ray))]
          [c (- (v-length-squared oc) (expt (sphere-radius sphere) 2))]
          [discriminant (- (* b/2 b/2) (* a c))])
         (if (< discriminant 0)
          -1
          ;else
          (/ (- (- b/2) (sqrt discriminant)) a))))

(struct hit-record (point normal t))
(define hittable-interface 
 (interface () hit))

(define sphere% 
  (class object%
   (hittable-interface)

   (init center radius)
   (define this-center center)
   (define this-radius radius)

   (super-new)

   (define/public (hit ray t-min t-max) 
    (letrec ([oc (v-sub (ray-origin ray) this-center)]
             [a (v-length-squared (ray-direction ray))]
             [b/2 (v-dot-prod oc (ray-direction ray))]
             [c (- (v-length-squared oc) (* this-radius this-radius))]
             [discriminant (- (* b/2 b/2) (* a c))])

          (if (< discriminant 0)
           '(#f)
          ;else
            (letrec ([dis-root (sqrt discriminant)]
                     [root (/ (- (- b/2) dis-root) a)])
                  
                    (cond [(or (< root t-min) (< t-max root) 
                            (let ([root (/ (+ (- b/2) dis-root) a)])
                              (if (or (< root t-min) (< t-max root))
                               '(#f)
                              ;else 
                               (list #t (hit-record (ray-at ray root) 
                                                    (v-div-scalar (v-sub (ray-at ray root) this-center) this-radius)
                                                    root)))))]
                        
                      [else (list #t (hit-record (ray-at ray root) 
                                      (v-div-scalar (v-sub (ray-at ray root) this-center) this-radius)
                                      root))]))))))) 
; ENOUGH
; I'm not going to implement this further for two reasons:
; 1. I don't know how to properly use Racket
; 2. I am following an imperative tutorial, and using a functional language
; 3. It does not make sense to write imperative code in a functional language and it is a pain    