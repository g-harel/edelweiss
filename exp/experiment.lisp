;; Genereates sequence [n..1].
(defun seq (count)
    (if (> count 0)
        (cons count (seq (- count 1)))))

(defun reduceRight (acc list handler)
    (if (> (length list) 0)
        (reduceRight
            (funcall handler acc (car list))
            (cdr list)
            handler)
        acc))

(defun adder0 (n)
    (reduceRight 0 (seq n) '+))

;;

(defun adder1 (n)
    (if (= n 0)
        0
        (+ n (adder1 (- n 1)))))

;;

(print (adder0 5))
(print (adder1 5))
