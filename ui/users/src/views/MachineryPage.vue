<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Refresh, User } from '@element-plus/icons-vue'
import AppFooter from './components/footer.vue'
import { useUserStore } from '../stores/user'
import { leaseApprovalApi, leaseProductApi } from '../services/api'
import type { LeaseApplicationRequest, LeaseApplication } from '../services/api'
import '../assets/icons/agri-icons.css'
import '../assets/icons/machinery-icons.css'
import MachineryCarousel from '../components/MachineryCarousel.vue'
import { carouselImages } from '../assets/images/machinery-images'

// 导入图片资源
import jiqi1 from '../assets/images/jiqi1.png';
import jiqi2 from '../assets/images/jiqi2.png';
import jiqi3 from '../assets/images/jiqi3.png';

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('equipment')
const loading = ref(false)
const submitting = ref(false)
const refreshing = ref(false)

// 当前订单信息
const orderInfo = ref({
  machineId: 'HN20240517',
  machineName: '马达5100联合收割机',
  machineType: '大型',
  power: '100马力',
  workWidth: '3.6米',
  fuelType: '柴油',
  rentalPrice: '600元/天',
  company: 'XX农业机械有限公司',
  pickupLocation: 'XX市XX区科学园惠农农机服务中心',
  returnLocation: 'XX市XX区科学园惠农农机服务中心',
  returnTime: '2024-06-20',
  estimatedFee: '1800',
  contactName: '',
  contactPhone: '',
  agreeTerms: false
})

// 农机类型数据
const machineryTypes = [
  { id: 1, name: '收割机', count: 5, icon: 'harvester' },
  { id: 2, name: '拖拉机', count: 2, icon: 'tractor' },
  { id: 3, name: '播种机', count: 2, icon: 'seeder' },
  { id: 4, name: '插秧机', count: 2, icon: 'transplanter' },
  { id: 5, name: '旋耕机', count: 2, icon: 'rotary-tiller' }
]

// 农机设备接口
interface MachineryItem {
  id: string
  name: string
  description: string
  category: string
  brand: string
  model: string
  power: string
  workWidth: string
  fuelType: string
  dailyPrice: number
  hourlyPrice: number
  available: boolean
  location: string
  contactPhone: string
  images: string[]
  specifications: Record<string, string>
}

// 农机设备列表
const machineryList = ref<MachineryItem[]>([])
const selectedCategory = ref('')

// 农机分类
const categories = [
  { value: '', label: '全部设备', icon: 'all' },
  { value: '收割机', label: '收割机', icon: 'harvester' },
  { value: '拖拉机', label: '拖拉机', icon: 'tractor' },
  { value: '播种机', label: '播种机', icon: 'seeder' },
  { value: '插秧机', label: '插秧机', icon: 'transplanter' },
  { value: '旋耕机', label: '旋耕机', icon: 'rotary' },
]

// 筛选后的农机列表
const filteredMachinery = computed(() => {
  if (!selectedCategory.value) {
    return machineryList.value
  }
  return machineryList.value.filter(item => item.category === selectedCategory.value)
})

// 格式化价格
const formatPrice = (price: number) => {
  return price.toLocaleString()
}

// 模拟农机数据
const mockMachineryData: MachineryItem[] = [
  {
    id: 'machinery_001',
    name: '雷沃谷神GE50收割机',
    description: '大型联合收割机，适合水稻、小麦收割，效率高，操作简便',
    category: '收割机',
    brand: '雷沃重工',
    model: 'GE50',
    power: '150马力',
    workWidth: '4.2米',
    fuelType: '柴油',
    dailyPrice: 800,
    hourlyPrice: 35,
    available: true,
    location: '泰安市农机服务中心',
    contactPhone: '18888888888',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '发动机功率': '150马力',
      '作业幅宽': '4.2米',
      '燃料类型': '柴油',
      '整机重量': '8.5吨',
      '作业效率': '20-30亩/小时'
    }
  },
  {
    id: 'machinery_002',
    name: '东方红1104拖拉机',
    description: '中大型轮式拖拉机，动力强劲，适合耕地、播种等多种作业',
    category: '拖拉机',
    brand: '一拖集团',
    model: '1104',
    power: '110马力',
    workWidth: '2.5米',
    fuelType: '柴油',
    dailyPrice: 500,
    hourlyPrice: 25,
    available: true,
    location: '洛阳市农机租赁点',
    contactPhone: '18777777777',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '发动机功率': '110马力',
      '驱动方式': '四轮驱动',
      '燃料类型': '柴油',
      '整机重量': '4.8吨',
      '适用作业': '耕地、播种、施肥'
    }
  },
  {
    id: 'machinery_003',
    name: '常发CF505插秧机',
    description: '水稻插秧专用机械，插秧速度快，秧苗成活率高',
    category: '插秧机',
    brand: '常发农机',
    model: 'CF505',
    power: '35马力',
    workWidth: '1.8米',
    fuelType: '汽油',
    dailyPrice: 400,
    hourlyPrice: 20,
    available: false,
    location: '常州市农机服务站',
    contactPhone: '18666666666',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '35马力',
      '插秧行数': '6行',
      '燃料类型': '汽油',
      '整机重量': '850公斤',
      '作业效率': '8-12亩/小时'
    }
  },
  {
    id: 'machinery_004',
    name: '德邦大为2BFX-12播种机',
    description: '小麦玉米通用播种机，播种均匀，深度可调',
    category: '播种机',
    brand: '德邦大为',
    model: '2BFX-12',
    power: '配套60-80马力拖拉机',
    workWidth: '3.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 300,
    hourlyPrice: 15,
    available: true,
    location: '德州市农机合作社',
    contactPhone: '18555555555',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '播种行数': '12行',
      '行距': '250mm',
      '作业幅宽': '3.0米',
      '整机重量': '1.2吨',
      '适用作物': '小麦、玉米、大豆'
    }
  },
  {
    id: 'machinery_005',
    name: '华德1GQN-200旋耕机',
    description: '多功能旋耕机，耕作深度可调，适合不同土壤类型',
    category: '旋耕机',
    brand: '华德农机',
    model: '1GQN-200',
    power: '配套40-60马力拖拉机',
    workWidth: '2.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 250,
    hourlyPrice: 12,
    available: true,
    location: '潍坊市农机服务中心',
    contactPhone: '18444444444',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '耕作深度': '15-25cm',
      '作业幅宽': '2.0米',
      '刀片数量': '32片',
      '整机重量': '680公斤',
      '适用土壤': '粘土、沙土、壤土'
    }
  },
  {
    id: 'machinery_006',
    name: '久保田DC70收割机',
    description: '进口联合收割机，性能稳定，收割质量高',
    category: '收割机',
    brand: '久保田',
    model: 'DC70',
    power: '70马力',
    workWidth: '2.8米',
    fuelType: '柴油',
    dailyPrice: 700,
    hourlyPrice: 30,
    available: true,
    location: '济南市农机租赁站',
    contactPhone: '18333333333',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '70马力',
      '作业幅宽': '2.8米',
      '燃料类型': '柴油',
      '整机重量': '3.8吨',
      '作业效率': '15-20亩/小时'
    }
  },
  {
    id: 'machinery_007',
    name: '约翰迪尔6B系列拖拉机',
    description: '国际知名品牌，高效可靠，适合大面积作业',
    category: '拖拉机',
    brand: '约翰迪尔',
    model: '6B-1454',
    power: '145马力',
    workWidth: '3.2米',
    fuelType: '柴油',
    dailyPrice: 650,
    hourlyPrice: 30,
    available: true,
    location: '青岛市农机中心',
    contactPhone: '18222222222',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '发动机功率': '145马力',
      '驱动方式': '四轮驱动',
      '燃料类型': '柴油',
      '整机重量': '6.2吨',
      '适用作业': '深耕、播种、收获'
    }
  },
  {
    id: 'machinery_008',
    name: '沃得锐龙4LZ-8A收割机',
    description: '国产优质收割机，适合多种作物收割，性价比高',
    category: '收割机',
    brand: '沃得农机',
    model: '4LZ-8A',
    power: '125马力',
    workWidth: '3.8米',
    fuelType: '柴油',
    dailyPrice: 720,
    hourlyPrice: 32,
    available: true,
    location: '南京市农机服务站',
    contactPhone: '18111111111',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '发动机功率': '125马力',
      '作业幅宽': '3.8米',
      '燃料类型': '柴油',
      '整机重量': '7.2吨',
      '作业效率': '18-25亩/小时'
    }
  },
  {
    id: 'machinery_009',
    name: '洋马VP8G插秧机',
    description: '日本进口高端插秧机，插秧精度高，作业效率佳',
    category: '插秧机',
    brand: '洋马农机',
    model: 'VP8G',
    power: '40马力',
    workWidth: '2.4米',
    fuelType: '汽油',
    dailyPrice: 480,
    hourlyPrice: 22,
    available: true,
    location: '杭州市农机租赁中心',
    contactPhone: '18999999999',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '40马力',
      '插秧行数': '8行',
      '燃料类型': '汽油',
      '整机重量': '1.1吨',
      '作业效率': '12-16亩/小时'
    }
  },
  {
    id: 'machinery_010',
    name: '大华宝来2BQX-200播种机',
    description: '气力式精密播种机，适合玉米、大豆等作物播种',
    category: '播种机',
    brand: '大华宝来',
    model: '2BQX-200',
    power: '配套80-120马力拖拉机',
    workWidth: '4.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 380,
    hourlyPrice: 18,
    available: true,
    location: '哈尔滨市农机合作社',
    contactPhone: '18188888888',
    images: ['./src/assets/images/jiqi1.png'],
    specifications: {
      '播种行数': '16行',
      '行距': '250mm',
      '作业幅宽': '4.0米',
      '整机重量': '1.8吨',
      '适用作物': '玉米、大豆、花生'
    }
  },
  {
    id: 'machinery_011',
    name: '金达威1GKN-300旋耕机',
    description: '重型旋耕机，适合坚硬土壤耕作，耐用性强',
    category: '旋耕机',
    brand: '金达威',
    model: '1GKN-300',
    power: '配套60-100马力拖拉机',
    workWidth: '3.0米',
    fuelType: '拖拉机带动',
    dailyPrice: 320,
    hourlyPrice: 16,
    available: true,
    location: '郑州市农机服务中心',
    contactPhone: '18177777777',
    images: ['./src/assets/images/jiqi2.png'],
    specifications: {
      '耕作深度': '18-30cm',
      '作业幅宽': '3.0米',
      '刀片数量': '48片',
      '整机重量': '950公斤',
      '适用土壤': '重粘土、硬质土壤'
    }
  },
  {
    id: 'machinery_012',
    name: '中联重科PL40水稻收割机',
    description: '专业水稻收割机，湿地作业能力强，收割损失率低',
    category: '收割机',
    brand: '中联重科',
    model: 'PL40',
    power: '90马力',
    workWidth: '2.2米',
    fuelType: '柴油',
    dailyPrice: 580,
    hourlyPrice: 26,
    available: false,
    location: '长沙市农机服务站',
    contactPhone: '18166666666',
    images: ['./src/assets/images/jiqi3.png'],
    specifications: {
      '发动机功率': '90马力',
      '作业幅宽': '2.2米',
      '燃料类型': '柴油',
      '整机重量': '4.5吨',
      '作业效率': '12-18亩/小时'
    }
  }
]

// 获取农机列表 - 使用真实的租赁产品API
const getMachineryList = async () => {
  loading.value = true
  try {
    const params = {
      page: 1,
      size: 20,
      status: 1, // 只获取可用的产品
      ...(selectedCategory.value && { type: selectedCategory.value })
    }
    
    const response = await leaseProductApi.getProducts(params)
    
    // 将租赁产品数据转换为农机数据格式
    machineryList.value = response.list.map(product => ({
      id: product.product_code,
      name: product.name,
      description: product.description,
      category: product.type,
      brand: product.brand,
      model: product.model,
      power: `${product.machinery}马力`, // 使用machinery字段作为动力
      workWidth: '详情咨询', // API中没有此字段，设置默认值
      fuelType: '柴油', // 设置默认值
      dailyPrice: product.daily_rate,
      hourlyPrice: Math.round(product.daily_rate / 8), // 按8小时计算时租
      available: product.available_count > 0,
      location: '服务中心', // 设置默认值
      contactPhone: '400-123-4567', // 设置默认值
      images: ['/src/assets/images/jiqi1.png'], // 设置默认图片
      specifications: {
        '产品编码': product.product_code,
        '品牌': product.brand,
        '型号': product.model,
        '日租金': `${product.daily_rate}元/天`,
        '押金': `${product.deposit}元`,
        '最长租期': `${product.max_duration}天`,
        '最短租期': `${product.min_duration}天`,
        '库存数量': `${product.inventory_count}台`,
        '可用数量': `${product.available_count}台`
      }
    }))
  } catch (error) {
    console.error('获取农机列表失败:', error)
    ElMessage.error('获取农机列表失败')
    // 如果API失败，使用模拟数据作为备用
    machineryList.value = mockMachineryData
  } finally {
    loading.value = false
  }
}

// 刷新农机列表
const onRefresh = async () => {
  refreshing.value = true
  try {
    await getMachineryList()
    ElMessage.success('刷新成功')
  } finally {
    refreshing.value = false
  }
}

// 查看农机详情
const viewMachineryDetail = (machinery: MachineryItem) => {
  // 跳转到农机详情页
  router.push(`/machinery/detail/${machinery.id}`)
}

// 检查库存可用性
const checkInventoryAvailability = async (productCode: string, startDate: string, endDate: string, quantity: number = 1) => {
  try {
    const response = await leaseProductApi.checkInventoryAvailability({
      product_code: productCode,
      quantity,
      start_date: startDate,
      end_date: endDate
    })
    return response
  } catch (error) {
    console.error('检查库存失败:', error)
    return { available: false, available_count: 0 }
  }
}

// 租赁农机
const rentMachinery = async (machinery: MachineryItem) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  if (!machinery.available) {
    ElMessage.warning('该设备暂不可租用')
    return
  }
  
  // 简单的库存检查（使用明天和后天作为示例）
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const dayAfterTomorrow = new Date()
  dayAfterTomorrow.setDate(dayAfterTomorrow.getDate() + 2)
  
  const startDate = tomorrow.toISOString().split('T')[0]
  const endDate = dayAfterTomorrow.toISOString().split('T')[0]
  
  // 检查库存
  const inventoryCheck = await checkInventoryAvailability(machinery.id, startDate, endDate)
  
  if (!inventoryCheck.available) {
    ElMessage.warning(`该设备在${startDate}至${endDate}期间暂无库存`)
    return
  }
  
  // 预填充表单数据
  formData.product_code = machinery.id
  formData.name = `${machinery.name}租赁申请`
  formData.type = machinery.category
  formData.machinery = machinery.name
  formData.daily_rate = machinery.dailyPrice
  formData.deposit = parseFloat(machinery.specifications['押金']?.replace('元', '') || '0')
  formData.start_date = startDate
  formData.end_date = endDate
  formData.duration = 1
  formData.total_amount = machinery.dailyPrice * 1
  formData.contact_phone = userStore.userInfo?.phone || ''
  
  // 切换到申请表单
  activeTab.value = 'apply'
  
  ElMessage.success(`已为您预填充${machinery.name}的租赁信息`)
}

// 联系租赁方
const contactRenter = (machinery: MachineryItem) => {
  const phone = machinery.contactPhone
  if (confirm(`确定要拨打 ${phone} 吗？`)) {
    window.location.href = `tel:${phone}`
  }
}

// 选择分类
const selectCategory = (category: string) => {
  selectedCategory.value = category
}

// 可用的租赁类型
const availableTypes = ref<string[]>([])

// 我的租赁申请
const myApplications = ref<LeaseApplication[]>([])

// 租赁类型配置
const typeConfigs = {
  '农机租赁': {
    icon: '🚜',
    description: '各类农业机械设备租赁服务',
    features: ['春耕必备', '高效作业', '专业维护'],
    examples: ['拖拉机', '收割机', '播种机', '施肥机']
  },
  '运输设备': {
    icon: '🚛',
    description: '农产品运输车辆租赁',
    features: ['货运便利', '冷链保鲜', '及时配送'],
    examples: ['货车', '冷藏车', '农用三轮车', '小型货车']
  },
  '加工设备': {
    icon: '⚙️',
    description: '农产品加工机械租赁',
    features: ['提升价值', '规模加工', '技术支持'],
    examples: ['打包机', '烘干机', '分拣机', '包装机']
  },
  '灌溉设备': {
    icon: '💧',
    description: '农田灌溉系统设备租赁',
    features: ['节水高效', '智能控制', '覆盖广泛'],
    examples: ['喷灌设备', '滴灌系统', '水泵', '管道']
  },
  '其他设备': {
    icon: '🔧',
    description: '其他农业相关设备租赁',
    features: ['种类丰富', '按需租赁', '灵活配置'],
    examples: ['发电机', '照明设备', '测量仪器', '安全设备']
  }
}

// 申请表单数据 - 更新字段名以匹配新API
const formData = reactive<LeaseApplicationRequest>({
  product_id: 1, // 默认产品ID，实际应该从产品选择中获取
  product_code: 'FARM_MACHINERY', // 默认产品编码
  name: '',
  type: '',
  machinery: '',
  start_date: '',
  end_date: '',
  duration: 0,
  daily_rate: 0,
  total_amount: 0,
  deposit: 0,
  delivery_address: '',
  contact_phone: '',
  purpose: ''
})

// 状态映射
const statusMap = {
  'pending': { text: '待审批', color: '#E6A23C', icon: '⏳' },
  'approved': { text: '已批准', color: '#67C23A', icon: '✅' },
  'rejected': { text: '已拒绝', color: '#F56C6C', icon: '❌' }
}

// 获取当前选择类型的配置
const currentTypeConfig = computed(() => {
  if (!formData.type) return null
  return typeConfigs[formData.type as keyof typeof typeConfigs]
})

// 格式化时间
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 格式化详细时间
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取状态信息
const getStatusInfo = (status: string) => {
  return statusMap[status as keyof typeof statusMap] || {
    text: status,
    color: '#909399',
    icon: '❓'
  }
}

// 计算租赁天数
const leaseDays = computed(() => {
  if (!formData.start_date || !formData.end_date) return 0
  const start = new Date(formData.start_date)
  const end = new Date(formData.end_date)
  const diffTime = Math.abs(end.getTime() - start.getTime())
  const days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  // 同步更新duration字段
  formData.duration = days
  return days
})

// 加载租赁类型 - 使用静态数据替代API调用
const loadLeaseTypes = async () => {
  try {
    // 新的API暂时没有getTypes方法，使用静态数据
    availableTypes.value = ['农机租赁', '运输设备', '加工设备', '灌溉设备', '其他设备']
  } catch (error: any) {
    console.error('加载租赁类型失败:', error)
    ElMessage.error('加载租赁类型失败')
  }
}

// 加载我的申请
const loadMyApplications = async () => {
  try {
    loading.value = true
    const response = await leaseApprovalApi.getMyApprovals({ page: 1, size: 10 })
    myApplications.value = response.list || []
  } catch (error: any) {
    console.error('加载我的申请失败:', error)
    ElMessage.error('加载我的申请失败')
  } finally {
    loading.value = false
  }
}

// 选择租赁类型
const selectType = (type: string) => {
  formData.type = type
  if (!formData.name) {
    formData.name = `${type}申请`
  }
  activeTab.value = 'apply'
}

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.type = ''
  formData.machinery = ''
  formData.start_date = ''
  formData.end_date = ''
  formData.duration = 0
  formData.daily_rate = 0
  formData.total_amount = 0
  formData.deposit = 0
  formData.delivery_address = ''
  formData.contact_phone = ''
  formData.purpose = ''
}

// 提交申请
const submitApplication = async () => {
  // 表单验证
  if (!formData.name.trim()) {
    ElMessage.error('请输入申请名称')
    return
  }
  if (!formData.type) {
    ElMessage.error('请选择租赁类型')
    return
  }
  if (!formData.start_date) {
    ElMessage.error('请选择开始时间')
    return
  }
  if (!formData.end_date) {
    ElMessage.error('请选择结束时间')
    return
  }
  if (new Date(formData.end_date) <= new Date(formData.start_date)) {
    ElMessage.error('结束时间必须晚于开始时间')
    return
  }
  if (!formData.purpose.trim() || formData.purpose.length < 10) {
    ElMessage.error('请输入详细的租赁用途说明，不少于10个字符')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认提交${formData.type}申请吗？租赁期间：${leaseDays.value}天`,
      '确认提交',
      {
        confirmButtonText: '确认提交',
        cancelButtonText: '再检查一下',
        type: 'warning'
      }
    )

    submitting.value = true
    const response = await leaseApprovalApi.create(formData)
    
    ElMessage.success('申请提交成功')
    
    // 重置表单
    resetForm()
    
    // 切换到我的申请标签页并刷新
    activeTab.value = 'my-applications'
    loadMyApplications()

  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('提交申请失败:', error)
      ElMessage.error(error.message || '提交申请失败')
    }
  } finally {
    submitting.value = false
  }
}



// 删除申请 - 新API使用cancel方法而不是delete
const deleteApplication = async (application: LeaseApplication) => {
  if (application.status !== 'pending') {
    ElMessage.warning('只能删除待审批的申请')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认删除申请"${application.name}"吗？删除后无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 新API使用cancel方法取消申请
    await leaseApprovalApi.cancel(application.application_id, '用户主动删除申请')
    ElMessage.success('删除成功')
    
    // 重新加载列表
    loadMyApplications()

  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('删除申请失败:', error)
      ElMessage.error(error.message || '删除申请失败')
    }
  }
}

// 查看申请详情
const viewApplicationDetail = (applicationId: string) => {
  router.push(`/lease/application/${applicationId}`)
}

// 标签页切换
const handleTabChange = (tabName: string) => {
  if (tabName === 'my-applications') {
    loadMyApplications()
  }
}

onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  // 初始化数据
  loadLeaseTypes()
  getMachineryList()
})
</script>

<template>
  <div class="machinery-container">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left">
        <el-icon @click="router.go(-1)"><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">农机租赁</div>
      <div class="nav-right">
        <el-icon @click="onRefresh" :class="{ 'is-loading': refreshing }">
          <Refresh />
        </el-icon>
      </div>
    </div>

    <div class="page-content">
      <!-- Tab切换 -->
      <div class="tab-container">
        <div class="tab-header">
          <div 
            class="tab-item"
            :class="{ active: activeTab === 'equipment' }"
            @click="activeTab = 'equipment'"
          >
            设备列表
          </div>
          <div 
            class="tab-item"
            :class="{ active: activeTab === 'apply' }"
            @click="activeTab = 'apply'"
          >
            租赁申请
          </div>
          <div 
            class="tab-item"
            :class="{ active: activeTab === 'my-applications' }"
            @click="activeTab = 'my-applications'"
          >
            我的申请
          </div>
        </div>
      </div>

      <!-- 设备列表Tab -->
      <div v-if="activeTab === 'equipment'" class="equipment-tab">
        <!-- 用户快捷操作 -->
        <div class="quick-actions" v-if="userStore.isLoggedIn">
        <div class="action-card primary" @click="router.push('/machinery/my-applications')">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="22" height="22" fill="white">
              <path d="M14,2H6C4.9,2,4,2.9,4,4v16c0,1.1,0.9,2,2,2h12c1.1,0,2-0.9,2-2V8L14,2z M16,18H8v-2h8V18z M16,14H8v-2h8V14z M13,9V3.5L18.5,9H13z"/>
            </svg>
          </div>
          <div class="card-content">
            <h3>我的申请</h3>
            <p>查看租赁申请进度</p>
          </div>
          <div class="card-arrow">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="white" opacity="0.8">
              <path d="M8.59,16.59L13.17,12L8.59,7.41L10,6l6,6l-6,6L8.59,16.59z"/>
            </svg>
          </div>
        </div>
        
        <div class="stats-row">
          <div class="stat-item">
            <div class="stat-value">2</div>
            <div class="stat-label">申请笔数</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M14,10H2V12H14V10M14,6H2V8H14V6M2,16H10V14H2V16M21.5,11.5L23,13L16,20L11.5,15.5L13,14L16,17L21.5,11.5Z"/>
              </svg>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-value">3台</div>
            <div class="stat-label">已租数量</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M18,18.5A1.5,1.5 0 0,1 16.5,17A1.5,1.5 0 0,1 18,15.5A1.5,1.5 0 0,1 19.5,17A1.5,1.5 0 0,1 18,18.5M19.5,9.5L21.46,12H17V9.5M6,18.5A1.5,1.5 0 0,1 4.5,17A1.5,1.5 0 0,1 6,15.5A1.5,1.5 0 0,1 7.5,17A1.5,1.5 0 0,1 6,18.5M20,8H17V4H3C1.89,4 1,4.89 1,6V17H3A3,3 0 0,0 6,20A3,3 0 0,0 9,17H15A3,3 0 0,0 18,20A3,3 0 0,0 21,17H23V12L20,8Z"/>
              </svg>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-value">4.8</div>
            <div class="stat-label">信用评分</div>
            <div class="stat-icon">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="#27ae60">
                <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- 登录提示 -->
      <div class="login-prompt" v-if="!userStore.isLoggedIn">
        <div class="prompt-content">
          <el-icon class="prompt-icon"><User /></el-icon>
          <p>登录后享受更多租赁服务</p>
          <el-button type="primary" @click="router.push('/login')">
            立即登录
          </el-button>
        </div>
      </div>

      <!-- 我的租赁入口 -->
      <div class="my-rental-card" v-if="userStore.isLoggedIn">
        <div class="rental-card-content" @click="router.push('/machinery/my-rentals')">
          <div class="rental-card-left">
            <div class="rental-card-icon">
              <svg viewBox="0 0 24 24" width="24" height="24" fill="white">
                <path d="M18,6H9V4H18M18,14v4H9v2H7v-2H3v-4H7v-2H9v2M7,14H5v2H7m11-2v6H9m9-10H9V6h9m0-4H9A2,2 0 0,0 7,4V12H9V8h9V20H9v2h11V4A2,2 0 0,0 18,2z"/>
              </svg>
            </div>
            <div class="rental-card-info">
              <h3 class="rental-card-title">我的租赁</h3>
              <p class="rental-card-desc">查看租赁记录及设备状态</p>
            </div>
          </div>
          <div class="rental-card-right">
            <span class="rental-card-count">3</span>
            <span class="rental-card-label">条记录</span>
            <svg viewBox="0 0 24 24" width="20" height="20" fill="#27ae60">
              <path d="M8.59,16.59L13.17,12L8.59,7.41L10,6l6,6l-6,6L8.59,16.59z"/>
            </svg>
          </div>
        </div>
      </div>

      <!-- 分类筛选 -->
      <div class="category-section">
        <div class="category-title">设备分类</div>
        <div class="category-tabs">
          <div
            v-for="category in categories"
            :key="category.value"
            class="category-tab"
            :class="{ active: selectedCategory === category.value }"
            @click="selectCategory(category.value)"
          >
            <div class="category-icon" :class="`icon-${category.icon}`"></div>
            <div class="category-name">{{ category.label }}</div>
          </div>
        </div>
      </div>

      <!-- 农机列表 -->
      <div class="machinery-section">
        <div class="section-header">
          <div class="section-title">可租设备</div>
          <div class="section-count">{{ filteredMachinery.length }}台设备</div>
        </div>

        <div v-loading="loading" class="machinery-list">
          <div
            v-for="machinery in filteredMachinery"
            :key="machinery.id"
            class="machinery-card"
          >
            <div class="card-header">
              <div class="machinery-info">
                <h3 class="machinery-name">{{ machinery.name }}</h3>
                <div class="machinery-meta">
                  <span class="category-tag">{{ machinery.category }}</span>
                  <span class="brand-tag">{{ machinery.brand }}</span>
                  <span class="availability-tag" :class="{ available: machinery.available, unavailable: !machinery.available }">
                    {{ machinery.available ? '可租用' : '已租出' }}
                  </span>
                </div>
                <p class="machinery-desc">{{ machinery.description }}</p>
              </div>
            </div>

            <div class="card-content">
              <div class="spec-grid">
                <div class="spec-item">
                  <span class="spec-label">动力</span>
                  <span class="spec-value">{{ machinery.power }}</span>
                </div>
                <div class="spec-item">
                  <span class="spec-label">幅宽</span>
                  <span class="spec-value">{{ machinery.workWidth }}</span>
                </div>
                <div class="spec-item">
                  <span class="spec-label">燃料</span>
                  <span class="spec-value">{{ machinery.fuelType }}</span>
                </div>
                <div class="spec-item">
                  <span class="spec-label">位置</span>
                  <span class="spec-value">{{ machinery.location }}</span>
                </div>
              </div>

              <div class="price-section">
                <div class="price-info">
                  <div class="daily-price">
                    <span class="price-label">日租</span>
                    <span class="price-value">¥{{ formatPrice(machinery.dailyPrice) }}</span>
                    <span class="price-unit">/天</span>
                  </div>
                  <div class="hourly-price">
                    <span class="price-label">时租</span>
                    <span class="price-value">¥{{ formatPrice(machinery.hourlyPrice) }}</span>
                    <span class="price-unit">/小时</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="card-actions">
              <button class="action-btn detail-btn" @click="viewMachineryDetail(machinery)">
                查看详情
              </button>
              <button
                class="action-btn contact-btn"
                @click="contactRenter(machinery)"
              >
                联系租赁
              </button>
              <button
                class="action-btn rent-btn"
                :class="{ disabled: !machinery.available }"
                @click="rentMachinery(machinery)"
                :disabled="!machinery.available"
              >
                {{ machinery.available ? '立即租赁' : '已租出' }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="filteredMachinery.length === 0 && !loading" class="empty-state">
          <div class="empty-icon">🚜</div>
          <div class="empty-text">暂无可租用设备</div>
          <div class="empty-desc">请稍后再试或联系客服</div>
        </div>
      </div>
    </div>
    </div>

    <!-- 租赁申请Tab -->
    <div v-if="activeTab === 'apply'" class="apply-tab">
      <div class="apply-form-container">
        <div class="form-header">
          <h3>租赁申请</h3>
          <p>请填写以下信息提交租赁申请</p>
        </div>

        <el-form :model="formData" label-width="100px" class="rental-form">
          <div class="form-section">
            <div class="section-title">基本信息</div>
            
            <el-form-item label="申请名称" required>
              <el-input v-model="formData.name" placeholder="请输入申请名称" />
            </el-form-item>
            
            <el-form-item label="设备类型" required>
              <el-select v-model="formData.type" placeholder="请选择设备类型" style="width: 100%">
                <el-option value="收割机" label="收割机" />
                <el-option value="拖拉机" label="拖拉机" />
                <el-option value="播种机" label="播种机" />
                <el-option value="插秧机" label="插秧机" />
                <el-option value="旋耕机" label="旋耕机" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="设备名称" required>
              <el-input v-model="formData.machinery" placeholder="请输入设备名称" />
            </el-form-item>
          </div>

          <div class="form-section">
            <div class="section-title">租期和费用</div>
            
            <el-form-item label="开始日期" required>
                             <el-date-picker
                 v-model="formData.start_date"
                 type="date"
                 placeholder="选择开始日期"
                 format="YYYY-MM-DD"
                 value-format="YYYY-MM-DD"
                 style="width: 100%"
                 :disabled-date="(time: Date) => time.getTime() < Date.now() - 8.64e7"
               />
            </el-form-item>
            
            <el-form-item label="结束日期" required>
                             <el-date-picker
                 v-model="formData.end_date"
                 type="date"
                 placeholder="选择结束日期"
                 format="YYYY-MM-DD"
                 value-format="YYYY-MM-DD"
                 style="width: 100%"
                 :disabled-date="(time: Date) => time.getTime() < Date.now() - 8.64e7"
               />
            </el-form-item>
            
            <el-form-item label="日租金">
              <el-input-number
                v-model="formData.daily_rate"
                :precision="2"
                :step="1"
                :min="0"
                style="width: 100%"
              />
              <div class="form-tip">租赁天数：{{ leaseDays }}天，总费用：¥{{ (formData.daily_rate * leaseDays).toFixed(2) }}</div>
            </el-form-item>
            
            <el-form-item label="押金">
              <el-input-number
                v-model="formData.deposit"
                :precision="2"
                :step="1"
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <div class="form-section">
            <div class="section-title">配送信息</div>
            
            <el-form-item label="配送地址" required>
              <el-input 
                v-model="formData.delivery_address" 
                type="textarea"
                :rows="2"
                placeholder="请输入详细的配送地址"
              />
            </el-form-item>
            
            <el-form-item label="联系电话" required>
              <el-input v-model="formData.contact_phone" placeholder="请输入联系电话" />
            </el-form-item>
          </div>

          <div class="form-section">
            <div class="section-title">租赁用途</div>
            
            <el-form-item label="使用用途" required>
              <el-input 
                v-model="formData.purpose" 
                type="textarea"
                :rows="3"
                placeholder="请详细说明设备使用用途，不少于10个字符"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
          </div>

          <div class="form-actions">
            <el-button size="large" @click="resetForm">重置表单</el-button>
            <el-button type="primary" size="large" @click="submitApplication" :loading="submitting">
              {{ submitting ? '提交中...' : '提交申请' }}
            </el-button>
          </div>
        </el-form>
      </div>
    </div>

    <!-- 我的申请Tab -->
    <div v-if="activeTab === 'my-applications'" class="my-applications-tab">
      <div class="applications-container">
        <div class="applications-header">
          <h3>我的租赁申请</h3>
          <el-button @click="loadMyApplications" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>

        <div v-loading="loading" class="applications-list">
          <div 
            v-for="app in myApplications" 
            :key="app.application_id"
            class="application-card"
          >
            <div class="app-header">
              <div class="app-info">
                <h4 class="app-name">{{ app.name }}</h4>
                <div class="app-meta">
                  <span class="app-type">{{ app.type }}</span>
                  <span class="app-id">ID: {{ app.application_id }}</span>
                </div>
              </div>
              <div 
                class="status-badge"
                :class="app.status"
              >
                {{ getStatusInfo(app.status).text }}
              </div>
            </div>

            <div class="app-content">
              <div class="app-details">
                <div class="detail-row">
                  <span class="label">设备名称:</span>
                  <span class="value">{{ app.machinery }}</span>
                </div>
                <div class="detail-row">
                  <span class="label">租赁期间:</span>
                  <span class="value">{{ app.start_date }} 至 {{ app.end_date }}</span>
                </div>
                <div class="detail-row">
                  <span class="label">日租金:</span>
                  <span class="value">¥{{ app.daily_rate }}</span>
                </div>
                <div class="detail-row">
                  <span class="label">总费用:</span>
                  <span class="value">¥{{ app.total_amount }}</span>
                </div>
              </div>
            </div>

            <div class="app-actions">
              <el-button size="small" @click="viewApplicationDetail(app.application_id)">
                查看详情
              </el-button>
              <el-button 
                v-if="app.status === 'pending'"
                size="small" 
                type="danger" 
                @click="deleteApplication(app)"
              >
                取消申请
              </el-button>
            </div>
          </div>

          <div v-if="myApplications.length === 0 && !loading" class="empty-applications">
            <div class="empty-icon">📋</div>
            <div class="empty-text">暂无申请记录</div>
            <div class="empty-desc">您还没有提交过租赁申请</div>
            <el-button type="primary" @click="activeTab = 'apply'">立即申请</el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.machinery-container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 80px;
}

.page-content {
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left, .nav-right {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #2c3e50;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.quick-actions {
  margin-bottom: 16px;
}

/* Tab样式 */
.tab-container {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tab-header {
  display: flex;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 4px;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 12px 16px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.3s ease;
  color: #666;
  font-weight: 500;
}

.tab-item.active {
  background: #27ae60;
  color: white;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.tab-item:hover:not(.active) {
  background: #e8f5e8;
  color: #27ae60;
}

/* 申请表单样式 */
.apply-tab {
  padding: 0;
}

.apply-form-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-header {
  text-align: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.form-header h3 {
  margin: 0 0 8px;
  font-size: 20px;
  color: #2c3e50;
}

.form-header p {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.rental-form .form-section {
  margin-bottom: 24px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.rental-form .section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #27ae60;
}

.rental-form .form-tip {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 4px;
}

.form-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.form-actions .el-button {
  min-width: 120px;
}

/* 我的申请样式 */
.my-applications-tab {
  padding: 0;
}

.applications-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.applications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.applications-header h3 {
  margin: 0;
  font-size: 18px;
  color: #2c3e50;
}

.applications-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  border: 1px solid #e1e8ed;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s ease;
}

.application-card:hover {
  border-color: #27ae60;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.1);
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.app-name {
  margin: 0 0 4px;
  font-size: 16px;
  color: #2c3e50;
}

.app-meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #7f8c8d;
}

.app-type {
  background: #e8f5e8;
  color: #27ae60;
  padding: 2px 8px;
  border-radius: 4px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.pending {
  background: #fff3cd;
  color: #856404;
}

.status-badge.approved {
  background: #d1ecf1;
  color: #0c5460;
}

.status-badge.rejected {
  background: #f8d7da;
  color: #721c24;
}

.status-badge.cancelled {
  background: #e2e3e5;
  color: #383d41;
}

.app-details {
  margin-bottom: 16px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.detail-row .label {
  color: #7f8c8d;
  font-weight: 500;
}

.detail-row .value {
  color: #2c3e50;
}

.app-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.empty-applications {
  text-align: center;
  padding: 40px 20px;
  color: #7f8c8d;
}

.empty-applications .empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-applications .empty-text {
  font-size: 18px;
  margin-bottom: 8px;
}

.empty-applications .empty-desc {
  font-size: 14px;
  margin-bottom: 20px;
}

.action-card {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 16px;
  padding: 22px;
  color: white;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 16px;
  box-shadow: 0 8px 16px rgba(39, 174, 96, 0.2);
  position: relative;
  overflow: hidden;
}

.action-card:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(45deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  z-index: 1;
}

.action-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 28px rgba(39, 174, 96, 0.3);
}

.card-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  position: relative;
  z-index: 2;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.card-content {
  flex: 1;
  position: relative;
  z-index: 2;
}

.card-content h3 {
  margin: 0 0 6px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
  font-weight: 400;
}

.card-arrow {
  opacity: 0.8;
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
}

.stats-row:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #27ae60, #2ecc71);
}

.stat-item {
  text-align: center;
  position: relative;
  padding: 8px 0;
}

.stat-item:not(:last-child):after {
  content: '';
  position: absolute;
  right: -8px;
  top: 20%;
  height: 60%;
  width: 1px;
  background: #f0f0f0;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #27ae60;
  margin-bottom: 8px;
  position: relative;
}

.stat-label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.stat-icon {
  position: absolute;
  top: 8px;
  right: 16px;
  opacity: 0.2;
}

/* 我的租赁入口 */
.my-rental-card {
  background: white;
  border-radius: 16px;
  margin-bottom: 16px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  position: relative;
}

.my-rental-card:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(to bottom, #27ae60, #2ecc71);
}

.rental-card-content {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.rental-card-content:hover {
  background-color: #f9f9f9;
}

.rental-card-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.rental-card-icon {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 10px rgba(39, 174, 96, 0.25);
}

.rental-card-info {
  display: flex;
  flex-direction: column;
}

.rental-card-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px;
}

.rental-card-desc {
  font-size: 13px;
  color: #666;
  margin: 0;
}

.rental-card-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rental-card-count {
  font-size: 18px;
  font-weight: 700;
  color: #27ae60;
}

.rental-card-label {
  font-size: 14px;
  color: #666;
}

/* 登录提示 */
.login-prompt {
  background: white;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.prompt-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.prompt-icon {
  font-size: 32px;
  color: #27ae60;
}

/* 分类筛选 */
.category-section {
  margin-bottom: 16px;
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.category-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
}

.category-tabs {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
}

.category-tab {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 15px 10px;
  border-radius: 8px;
  background: #f8f9fa;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.category-tab:hover {
  background: #e9ecef;
  transform: translateY(-2px);
}

.category-tab.active {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  border-color: #27ae60;
  box-shadow: 0 4px 15px rgba(39, 174, 96, 0.3);
}

.category-icon {
  width: 24px;
  height: 24px;
  margin-bottom: 8px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
  text-align: center;
}

/* 农机列表 */
.machinery-section {
  margin-bottom: 16px;
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.section-count {
  font-size: 14px;
  color: #666;
}

.machinery-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.machinery-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
}

.machinery-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.card-header {
  padding: 20px 20px 0;
}

.machinery-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 10px;
}

.machinery-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.category-tag, .brand-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.category-tag {
  background: #e3f2fd;
  color: #1976d2;
}

.brand-tag {
  background: #f3e5f5;
  color: #7b1fa2;
}

.availability-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.availability-tag.available {
  background: #e8f5e8;
  color: #2e7d32;
}

.availability-tag.unavailable {
  background: #ffebee;
  color: #c62828;
}

.machinery-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0;
}

.card-content {
  padding: 15px 20px;
}

.spec-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
  margin-bottom: 15px;
}

.spec-item {
  text-align: center;
}

.spec-label {
  display: block;
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.spec-value {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.price-section {
  border-top: 1px solid #eee;
  padding-top: 15px;
}

.price-info {
  display: flex;
  justify-content: space-around;
  text-align: center;
}

.daily-price, .hourly-price {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.price-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.price-value {
  font-size: 18px;
  font-weight: 600;
  color: #27ae60;
}

.price-unit {
  font-size: 12px;
  color: #999;
}

.card-actions {
  display: flex;
  gap: 10px;
  padding: 15px 20px 20px;
}

.action-btn {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.detail-btn {
  background: #f8f9fa;
  color: #333;
}

.detail-btn:hover {
  background: #e9ecef;
}

.contact-btn {
  background: #fff3e0;
  color: #f57c00;
}

.contact-btn:hover {
  background: #ffe0b2;
}

.rent-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
}

.rent-btn:hover:not(.disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(39, 174, 96, 0.3);
}

.rent-btn.disabled {
  background: #e0e0e0;
  color: #999;
  cursor: not-allowed;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #666;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .machinery-container {
    padding-bottom: 70px;
  }
  
  .page-content {
    padding: 12px;
  }
  
  .action-card {
    padding: 18px;
  }
  
  .card-icon {
    width: 42px;
    height: 42px;
  }
  
  .card-content h3 {
    font-size: 18px;
  }
  
  .stats-row {
    padding: 16px 10px;
  }
  
  .stat-value {
    font-size: 22px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .stat-icon {
    display: none;
  }
  
  .category-tabs {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .spec-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .card-actions {
    flex-direction: column;
  }
  
  .action-btn {
    flex: none;
  }
}
</style> 