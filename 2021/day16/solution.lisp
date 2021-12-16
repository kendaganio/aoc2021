(require :asdf)

; set-up phase
(defvar htob (make-hash-table :test 'equal))
(defvar v-total 0)
(loop for ch across "0123456789ABCDEF" do 
      (setf (gethash (string ch) htob) 
            (format nil "~v,'0B" 4 (parse-integer (string ch) :radix 16))))

(defun str-to-binary (str &aux out)
  (loop for ch across str do
        (setf out (concatenate 'string out (gethash (string ch) htob))))
  (identity out))

(defun parse-literal (binary-str &aux current remaining done (parsed ""))
  (setf remaining binary-str)
  (loop while (not done) do
        (setf current (subseq remaining 0 5))
        (setf parsed (format nil "~a~a" parsed (subseq current 1)))
        (if (equal (subseq current 0 1) "0")
          (setf done 't))
        (setf remaining (subseq remaining 5)))
  (cons (parse-integer parsed :radix 2) remaining))

(defun parse-operator (binary-str tid &aux res out ltid bit-len sub-len og-len remaining done)
  (setf remaining (subseq binary-str 1))
  (setf ltid (subseq binary-str 0 1))
  (if (equal ltid "0")
    (setf bit-len 15)
    (setf bit-len 11))

  (setf sub-len (parse-integer (subseq remaining 0 bit-len) :radix 2))
  (setf og-len (length (subseq remaining bit-len)))
  (setf remaining (subseq remaining bit-len))


  (loop while (not done) do
        (setf res (parse-bits remaining))
        (setf out (cons (car res) out))
        (setf remaining (cdr res))

        (if (equal ltid "0")
          (setf done (= og-len (+ sub-len (length remaining)))))

        (if (equal ltid "1")
          (setf done (equal (length out) sub-len)))) 

  (perform-operation (cons out remaining) tid))

(defun parse-bits (binary-str &aux v tid payload parsed-value)
  (setf v (format nil "~a" (parse-integer (subseq binary-str 0 3) :radix 2)))
  (setf v-total (+ v-total (parse-integer v)))
  (setf tid (format nil "~a" (parse-integer (subseq binary-str 3 6) :radix 2)))
  (setf payload (subseq binary-str 6))

  (if (equal tid "4") (parse-literal payload) (parse-operator payload tid)))

(defun perform-operation (res tid &aux out)
  (setf out (car res))
  (if (equal tid "0")
    (setf out (reduce #'+ out)))
  (if (equal tid "1")
    (setf out (reduce #'* out)))
  (if (equal tid "2")
    (setf out (apply #'min out)))
  (if (equal tid "3")
    (setf out (apply #'max out)))
  (if (equal tid "5")
    (if (< (car out) (cadr out))
      (setf out 1)
      (setf out 0)))
  (if (equal tid "6")
    (if (> (car out) (cadr out))
      (setf out 1)
      (setf out 0)))
  (if (equal tid "7")
    (if (= (car out) (cadr out))
      (setf out 1)
      (setf out 0)))
  (cons out (cdr res)))

(defun solve ()
  (format t "[PART2]: ~a~%" (car (parse-bits (str-to-binary "020D708041258C0B4C683E61F674A1401595CC3DE669AC4FB7BEFEE840182CDF033401296F44367F938371802D2CC9801A980021304609C431007239C2C860400F7C36B005E446A44662A2805925FF96CBCE0033C5736D13D9CFCDC001C89BF57505799C0D1802D2639801A900021105A3A43C1007A1EC368A72D86130057401782F25B9054B94B003013EDF34133218A00D4A6F1985624B331FE359C354F7EB64A8524027D4DEB785CA00D540010D8E9132270803F1CA1D416200FDAC01697DCEB43D9DC5F6B7239CCA7557200986C013912598FF0BE4DFCC012C0091E7EFFA6E44123CE74624FBA01001328C01C8FF06E0A9803D1FA3343E3007A1641684C600B47DE009024ED7DD9564ED7DD940C017A00AF26654F76B5C62C65295B1B4ED8C1804DD979E2B13A97029CFCB3F1F96F28CE43318560F8400E2CAA5D80270FA1C90099D3D41BE00DD00010B893132108002131662342D91AFCA6330001073EA2E0054BC098804B5C00CC667B79727FF646267FA9E3971C96E71E8C00D911A9C738EC401A6CBEA33BC09B8015697BB7CD746E4A9FD4BB5613004BC01598EEE96EF755149B9A049D80480230C0041E514A51467D226E692801F049F73287F7AC29CB453E4B1FDE1F624100203368B3670200C46E93D13CAD11A6673B63A42600C00021119E304271006A30C3B844200E45F8A306C8037C9CA6FF850B004A459672B5C4E66A80090CC4F31E1D80193E60068801EC056498012804C58011BEC0414A00EF46005880162006800A3460073007B620070801E801073002B2C0055CEE9BC801DC9F5B913587D2C90600E4D93CE1A4DB51007E7399B066802339EEC65F519CF7632FAB900A45398C4A45B401AB8803506A2E4300004262AC13866401434D984CA4490ACA81CC0FB008B93764F9A8AE4F7ABED6B293330D46B7969998021C9EEF67C97BAC122822017C1C9FA0745B930D9C480"))))
  (format t "[PART1]: ~a~%" v-total))

(time (solve))
