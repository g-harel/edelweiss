;; Genereates sequence [n..1].
(defun seq (count)
    (if (> count 0)
        (cons count (seq (- count 1)))))

;; Reduces list items into accumulator value using handler.
(defun reduceRight (acc list handler)
    (if (> (length list) 0)
        (reduceRight
            (funcall handler acc (car list))
            (cdr list)
            handler)
        acc))

;; Splits string into list of single-character strings.
(defun split (str)
    (loop for i from 0 to (- (length str) 1)
        collect (string (aref str i))))

;;
;;

(defun adder0 (n)
    (reduceRight 0 (seq n) '+))

;;

(defun adder1 (n)
    (if (= n 0)
        0
        (+ n (adder1 (- n 1)))))

;;

(print
    (adder0 5))
(print
    (adder1 5))

;;
;;

(defun counter0 (str char)
    (reduceRight
        0
        (split str)
        (lambda (acc currChar)
            (if (string-equal currChar char)
                (+ 1 acc)
                acc))))

(print
    (counter0 "abaasdass" "a"))
