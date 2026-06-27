export const featuredArticle = {
  id: 1,
  title: '把个人博客做成长期作品，而不只是练手项目',
  cover:
    'https://images.unsplash.com/photo-1499750310107-5fef28a66643?auto=format&fit=crop&w=1200&q=80',
  excerpt:
    '从内容结构、视觉表达到工程落地，博客真正的价值在于它能够持续记录成长轨迹，并逐步形成自己的技术品牌。',
  date: '2026-04-28',
  readTime: '6 分钟',
  category: '思考',
}

export const articles = [
  {
    id: 1,
    title: '把个人博客做成长期作品，而不只是练手项目',
    excerpt:
      '一个成熟的博客，不只是文章容器，而是一套完整的个人表达系统。',
    category: '思考',
    date: '2026-04-28',
    readTime: '6 分钟',
    tags: ['博客设计', '成长记录', '项目规划'],
    content: [
      '当我们开始写博客时，最容易关注的是技术栈和界面是否足够酷炫，但真正让博客有生命力的，其实是持续输出与内容组织能力。',
      '一个好的个人博客应该回答三个问题：你是谁、你擅长什么、你正在思考什么。首页负责建立第一印象，文章列表负责沉淀内容，关于页负责补充人物信息，后台则服务于未来的写作管理。',
      '因此，这套前端页面在视觉上采用了简洁、留白和卡片式布局，在结构上保留了文章、分类、标签、作者介绍和控制台入口，为你后续继续扩展后端接口预留空间。',
    ],
  },
  {
    id: 2,
    title: '前后端分离项目里，先做页面还是先补接口？',
    excerpt:
      '如果接口还没完善，先搭前端骨架往往更高效，因为它能帮助你提前明确数据结构和交互流程。',
    category: '开发',
    date: '2026-04-20',
    readTime: '5 分钟',
    tags: ['前后端分离', '接口设计', '开发流程'],
    content: [
      '很多个人项目在开发中期会停滞，原因通常不是不会写代码，而是缺少阶段性成果。先把前端页面搭起来，能让项目尽快具备可视化成果。',
      '当页面跑起来以后，哪些数据是列表、哪些数据要分页、哪些字段需要高亮，都会变得更加清晰。这样反过来会帮助你更好地定义后端 DTO 与接口返回格式。',
      '你现在的项目就很适合这个策略：登录先接真实接口，其余页面先用优质静态数据承载，等后端补完后再逐步替换。',
    ],
  },
  {
    id: 3,
    title: '如何让博客首页看起来更高级',
    excerpt:
      '高级感往往来自节奏、层次和留白，而不是堆砌过多的视觉元素。',
    category: '设计',
    date: '2026-04-14',
    readTime: '4 分钟',
    tags: ['UI', '视觉设计', '首页布局'],
    content: [
      '首页最重要的是首屏表达。一个清晰的标题、一段简洁的介绍、一组核心数据，再配合层次分明的色彩，就能建立相当不错的视觉质感。',
      '在这套设计中，我用了大面积渐变背景、柔和阴影和半透明卡片，既保留了现代风格，也不会显得过于花哨。',
      '同时，文章卡片和信息模块统一了圆角、间距与文字层级，让整个站点具备连续的视觉语言。',
    ],
  },
]

export const categories = ['全部', '思考', '开发', '设计']

export const stats = [
  { label: '已发布文章', value: '24+' },
  { label: '累计阅读', value: '18.6k' },
  { label: '专题分类', value: '8' },
  { label: '写作年限', value: '3 年' },
]

export const dashboardCards = [
  { label: '草稿箱', value: '12', trend: '+3 本周' },
  { label: '待回复评论', value: '18', trend: '+5 今日' },
  { label: '新增访客', value: '1,284', trend: '+12.8%' },
  { label: '收藏次数', value: '326', trend: '+42 本月' },
]

export const profile = {
  name: 'Jun',
  title: '后端开发者 / 博客作者',
  intro:
    '专注 Go 后端开发，也热爱把技术过程整理成可阅读、可分享的内容。这个博客会记录开发实践、项目复盘和一些关于成长的思考。',
  bio: '一个喜欢用代码解决问题的人。相信技术写作是思考的延伸，每一篇文章都是一次深度复盘。从后端架构到工程实践，记录成长路上的真实经验。',
  social: {
    github: 'https://github.com/super164',
    email: '2696057150@qq.com',
  },
  skills: [
    { name: 'Go', category: '语言' },
    { name: 'Java', category: '语言' },
    { name: 'Vue.js', category: '前端' },
    { name: 'React', category: '前端' },
    { name: 'MySQL', category: '数据库' },
    { name: 'Redis', category: '数据库' },
    { name: 'Docker', category: '工具' },
    { name: 'Git', category: '工具' },
    { name: 'Linux', category: '工具' },
    { name: 'Nginx', category: '工具' },
  ],
  timeline: [
      {
          year: '2026',
          title: '开始汇总个人技术写作',
          description: '将开发中的思考和实践整理成文章，汇总，建立个人博客。',
      },
      {
      year: '2025',
      title: '深入 Go 后端开发',
      description: '系统学习 Go 语言和后端架构设计，参与多个项目开发。',
    },
    {
      year: '2024',
      title: '系统性学习一门语言',
      description: '系统学习 JAVA语言',
    },
    {
      year: '2023',
      title: '踏入编程世界',
      description: '写下第一行代码，开启技术探索之旅。',
    },
  ],
}
