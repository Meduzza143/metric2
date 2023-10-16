package agent

// func TestSendData(t *testing.T) {

// 	tests := []struct { // добавляем слайс тестов
// 		name      string
// 		value     string
// 		valueType string
// 		want      int
// 	}{
// 		{
// 			name:      "test_1",
// 			value:     "55.12312",
// 			valueType: "gauge",
// 			want:      200,
// 		},
// 		{
// 			name:      "test_2",
// 			value:     "not_a_value",
// 			valueType: "gauge",
// 			want:      400,
// 		},
// 		{
// 			name:      "test_3",
// 			value:     "12",
// 			valueType: "not_a_type",
// 			want:      400,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			if resp := SendData(test.value, test.name, test.valueType); resp != test.want {
// 				t.Errorf("error has been accured: got %d (expected %d)", resp, test.want)
// 			}
// 		})
// 	}
// }
