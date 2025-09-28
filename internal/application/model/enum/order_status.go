package enum

import "errors"

// INITIATE：已发起
// SUCCESS：支付成功
// AFTERSALE：售后处理
// REFUND：已退款
// NOTPAY：未支付
// CLOSED：已关闭
// PAYERROR：支付失败(其他原因，如银行返回失败)
var ORDER_STATUS = [...]string{"", "INITIATE", "SUCCESS", "AFTERSALE", "REFUND", "NOTPAY", "CLOSED", "PAYERROR"}

type OrderStatus uint

// JsonDate反序列化
func (t *OrderStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range ORDER_STATUS {
		if status == value {
			*t = OrderStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t OrderStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ORDER_STATUS[t] + "\""), nil
}

func OrderStatusValue(value string) (OrderStatus, error) {
	for i, status := range ORDER_STATUS {
		if status == value {
			return OrderStatus(i), nil
		}
	}
	return OrderStatus(0), errors.New("未找到状态码：" + value)
}
