// dataProcessor.worker.js
// 工作线程中没有 window 对象，全局是 self

// 接收主线程发送的数据
self.onmessage = (e) => {
    const data = e.data; // 接收 10 万条数据
  
    // 耗时计算：统计总和、平均值、最大值
    let sum = 0;
    let max = 0;
    const len = data.length;
  
    for (let i = 0; i < len; i++) {
      const num = data[i];
      sum += num;
      if (num > max) max = num;
    }
  
    const avg = sum / len;
  
    // 将结果发送回主线程
    self.postMessage({ sum, avg, max });
  };